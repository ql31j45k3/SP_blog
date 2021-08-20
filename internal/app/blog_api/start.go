package blog_api

import (
	"fmt"

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

	container := buildContainer()

	// 調用其他函式，函式參數容器會依照 Provide 提供後自行匹配
	container.Invoke(article.RegisterRouter)
	container.Invoke(author.RegisterRouter)

	container.Invoke(func(r *gin.Engine) {
		utilsDriver.StartGin(r)
	})
}

// buildContainer 建立 DI 容器，提供各個函式的 input 參數
func buildContainer() *dig.Container {
	container := dig.New()
	provideFunc := containerProvide{}

	container.Provide(provideFunc.gin)
	container.Provide(provideFunc.gorm)
	container.Provide(provideFunc.translator)

	return container
}

type containerProvide struct {
	_ struct{}
}

// gin 建立 gin Engine，設定 middleware
func (cp *containerProvide) gin() *gin.Engine {
	return gin.Default()
}

// gorm 建立 gorm.DB 設定，初始化 session 並無實際連線
func (cp *containerProvide) gorm() (*gorm.DB, error) {
	return utilsDriver.NewMysql(configs.Gorm.GetHost(), configs.Gorm.GetUsername(), configs.Gorm.GetPassword(),
		configs.Gorm.GetDBName(), configs.Gorm.GetPort(), configs.Gorm.GetLogMode())
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
