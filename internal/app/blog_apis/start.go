package blog_apis

import (
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/configs"
	"go.uber.org/dig"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ql31j45k3/SP_blog/internal/binder"
	"github.com/ql31j45k3/SP_blog/internal/modules/article"
)

func Start() {
	container := buildContainer()

	container.Invoke(article.Router)
	container.Invoke(func(r *gin.Engine) {
		r.Run(configs.ConfigHost.GetSPBlogApisHost())
	})
}

func buildContainer() *dig.Container {
	container := binder.Container

	container.Provide(func() *gin.Engine {
		return gin.Default()
	})

	container.Provide(func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(configs.ConfigDB.GetDSN()), &gorm.Config{})
	})

	return container
}
