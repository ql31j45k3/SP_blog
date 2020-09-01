package article

import (
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils"
	"gorm.io/gorm"
)

func Router(r *gin.Engine, db *gorm.DB) {
	r.GET("/article/:id", func(c *gin.Context) {
		useCase := newUseCaseArticle(c, db)
		result, err := useCase.GetID()
		if err != nil {
			return
		}

		utils.ResponeOK(c, result)
	})
}
