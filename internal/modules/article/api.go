package article

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Router(r *gin.Engine, db *gorm.DB) {
	r.POST("/article", func(c *gin.Context) {
		useCase := newUseCaseArticle(c, db)
		result, err := useCase.Post()
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, result)
	})

	r.GET("/article/:id", func(c *gin.Context) {
		useCase := newUseCaseArticle(c, db)
		result, err := useCase.GetID()
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
