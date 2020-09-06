package blog_apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"
	"github.com/ql31j45k3/SP_blog/configs"
	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/ql31j45k3/SP_blog/internal/binder"
	"github.com/ql31j45k3/SP_blog/internal/modules/article"
	validatorFunc "github.com/ql31j45k3/SP_blog/internal/utils/validator"
)

// Start 控制服務流程、呼叫的依賴性
func Start() {
	container := buildContainer()

	// 調用其他函式，函式參數容器會依照 Provide 提供後自行匹配
	container.Invoke(article.SetupRouter)
	container.Invoke(func(r *gin.Engine) {
		// 控制調試日誌 log
		gin.SetMode(configs.ConfigGin.GetMode())

		r.Run(configs.ConfigHost.GetSPBlogApisHost())
	})
}

// buildContainer 建立 DI 容器，提供各個函式的 input 參數
func buildContainer() *dig.Container {
	container := binder.Container

	// 建立 gin Engine，設定 middleware
	container.Provide(func() *gin.Engine {
		return gin.Default()
	})

	// 建立 gorm.DB 設定，初始化 session 並無實際連線
	container.Provide(func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(configs.ConfigDB.GetDSN()), &gorm.Config{
			Logger: logger.Default.LogMode(configs.ConfigGorm.GetLogMode()),
		})
	})

	// 建立 Translator 設定翻譯語言類型、可自行擴充驗證函式與翻譯訊息 func
	container.Provide(func() ut.Translator {
		locale := configs.ConfigValidator.GetLocale()
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
				validatorFunc.ArticleStatusValidator); err != nil {
				panic(err)
			}

			// 根據提供的標記註冊翻譯
			v.RegisterTranslation(validatorFunc.ArticleStatusTag, trans,
				validatorFunc.ArticleStatusTranslations, validatorFunc.ArticleStatusTranslation)
		}

		return trans
	})

	return container
}
