package blog_api

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/ql31j45k3/SP_blog/configs"
	"github.com/ql31j45k3/SP_blog/internal/modules/article"
	"github.com/ql31j45k3/SP_blog/internal/modules/author"
	validatorFunc "github.com/ql31j45k3/SP_blog/internal/utils/validator"
	"go.uber.org/dig"
	"gorm.io/gorm"

	utilsDriver "github.com/ql31j45k3/SP_blog/internal/utils/driver"

	_ "net/http/pprof"
)

// Start 控制服務流程、呼叫的依賴性
func Start() {
	// 開始讀取設定檔，順序上必須為容器之前，執行容器內有需要設定檔 struct 取得參數
	if err := configs.Start(""); err != nil {
		panic(fmt.Errorf("start - configs.Start - %w", err))
	}

	// 順序必須在 configs 之後，需取得 設定參數
	if configs.IsPrintVersion() {
		return
	}

	// 開始讀取翻譯檔案，順序上必須在容器前執行
	validatorFunc.Start()

	utilsDriver.SetLogEnv()
	configs.SetReloadFunc(utilsDriver.ReloadSetLogLevel)

	go func() {
		if configs.Env.GetPPROFBlockStatus() {
			runtime.SetBlockProfileRate(configs.Env.GetPPROFBlockRate())
		}

		if configs.Env.GetPPROFMutexStatus() {
			runtime.SetMutexProfileFraction(configs.Env.GetPPROFMutexRate())
		}

		if configs.Env.GetPPROFStatus() {
			_ = http.ListenAndServe(configs.Host.GetPPROFAPIHost(), nil)
		}
	}()

	_, cancelCtxStopNotify := context.WithCancel(context.Background())
	// 注意: cancelCtx 底層保證多個調用，只會執行一次
	defer cancelCtxStopNotify()

	stopJobFunc := stopJob{}

	container, err := buildContainer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - buildContainer")
		return
	}

	// 調用其他函式，函式參數容器會依照 Provide 提供後自行匹配
	if err := container.Invoke(article.RegisterRouter); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - container.Invoke(article.RegisterRouter)")
		return
	}

	if err := container.Invoke(author.RegisterRouter); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - container.Invoke(author.RegisterRouter)")
		return
	}

	if err := container.Invoke(func(r *gin.Engine) {
		utilsDriver.StartGin(cancelCtxStopNotify, stopJobFunc.stop, r)
	}); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Start - utilsDriver.StartGin")
		return
	}
}

// stopJob 為避免其它 package 需 import 此包 package，故用傳遞 func 方式提供功能給其它模組使用，
// 依賴關係都是 start.go 單向 import 其它 package 包功能
type stopJob struct {
	sync.Mutex
	stopFunctions []func()
}

func (s *stopJob) stop() context.Context {
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func(s *stopJob, cancelCtx context.CancelFunc) {
		s.Lock()
		defer s.Unlock()

		defer cancelCtx()

		for _, f := range s.stopFunctions {
			f()
		}
	}(s, cancelCtx)

	return ctx
}

func (s *stopJob) add(f func()) {
	s.Lock()
	defer s.Unlock()

	s.stopFunctions = append(s.stopFunctions, f)
}

// buildContainer 建立 DI 容器，提供各個函式的 input 參數
func buildContainer() (*dig.Container, error) {
	container := dig.New()
	provideFunc := containerProvide{}

	if err := container.Provide(provideFunc.gin); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.gin) - %w", err)
	}

	if err := container.Provide(provideFunc.mysqlMaster, dig.Name("dbM")); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.mysqlMaster) - %w", err)
	}

	if err := container.Provide(provideFunc.translator); err != nil {
		return nil, fmt.Errorf("container.Provide(provideFunc.translator) - %w", err)
	}

	return container, nil
}

type containerProvide struct {
}

// gin 建立 gin Engine，設定 middleware
func (cp *containerProvide) gin() *gin.Engine {
	return gin.Default()
}

// gorm 建立 gorm.DB 設定，初始化 session 並無實際連線
func (cp *containerProvide) mysqlMaster() (*gorm.DB, error) {
	return utilsDriver.NewMysql(configs.Gorm.GetMasterHost(), configs.Gorm.GetMasterUsername(), configs.Gorm.GetMasterPassword(),
		configs.Gorm.GetMasterDBName(), configs.Gorm.GetMasterPort(), configs.Gorm.GetLogMode(),
		configs.Gorm.GetMasterMaxIdle(), configs.Gorm.GetMasterMaxOpen(), configs.Gorm.GetMasterMaxLifetime())
}

// translator 建立 Translator 設定翻譯語言類型、可自行擴充驗證函式與翻譯訊息 func
func (cp *containerProvide) translator() ut.Translator {
	locale := configs.Validator.GetLocale()
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator(locale)
	// 設定語言地區
	validatorFunc.SetLocale(locale)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 註冊翻譯器
		err := zhTranslations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			panic(err)
		}

		// 註冊一個函式，獲取struct tag裡自定義的label作為欄位名
		v.RegisterTagNameFunc(validatorFunc.RegisterTagNameFunc)

		// 註冊自定義函式
		if err := v.RegisterValidation(validatorFunc.ArticleStatusTag,
			validatorFunc.ArticleStatusFunc.Validator); err != nil {
			panic(err)
		}

		// 根據提供的標記註冊翻譯
		v.RegisterTranslation(validatorFunc.ArticleStatusTag, trans,
			validatorFunc.ArticleStatusFunc.Translations, validatorFunc.ArticleStatusFunc.Translation)

		// 註冊自定義函式
		if err := v.RegisterValidation(validatorFunc.AuthorStatusTag,
			validatorFunc.AuthorStatusFunc.Validator); err != nil {
			panic(err)
		}

		// 根據提供的標記註冊翻譯
		v.RegisterTranslation(validatorFunc.AuthorStatusTag, trans,
			validatorFunc.AuthorStatusFunc.Translations, validatorFunc.AuthorStatusFunc.Translation)
	}

	return trans
}
