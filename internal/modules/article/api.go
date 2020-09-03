package article

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	routerGroup := r.Group("/v1/article")

	routerGroup.POST("", func(c *gin.Context) {
		useCase := newUseCaseArticle(c, db)
		result, err := useCase.Post()
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, result)
	})

	routerGroup.GET("/:id", func(c *gin.Context) {
		useCase := newUseCaseArticle(c, db)
		result, err := useCase.GetID()
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
