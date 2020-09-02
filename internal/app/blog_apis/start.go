package blog_apis

import (
	_ "github.com/spf13/viper"

	"github.com/gin-gonic/gin"
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
		r.Run(":8080")
	})
}

const (
	dsn = "sp_blog_apis:sp_blog_apis@tcp(localhost:3306)/sp_blog?charset=utf8&parseTime=True&loc=Local"
)

func buildContainer() *dig.Container {
	container := binder.Container

	container.Provide(func() *gin.Engine {
		return gin.Default()
	})

	container.Provide(func() (*gorm.DB, error){
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})

	return container
}
