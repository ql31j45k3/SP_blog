package article

import (
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils"
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

		utils.ResponeOK(c, result)
	})
}
