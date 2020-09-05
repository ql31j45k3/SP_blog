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

func Start() {
	container := buildContainer()

	container.Invoke(article.SetupRouter)
	container.Invoke(func(r *gin.Engine) {
		gin.SetMode(configs.ConfigGin.GetMode())

		r.Run(configs.ConfigHost.GetSPBlogApisHost())
	})
}

func buildContainer() *dig.Container {
	container := binder.Container

	container.Provide(func() *gin.Engine {
		return gin.Default()
	})

	container.Provide(func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(configs.ConfigDB.GetDSN()), &gorm.Config{
			Logger: logger.Default.LogMode(configs.ConfigGorm.GetLogMode()),
		})
	})

	container.Provide(func() ut.Translator {
		uni := ut.New(zh.New())
		trans, _ := uni.GetTranslator("zh")

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
				validatorFunc.StatusValidator); err != nil {
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
