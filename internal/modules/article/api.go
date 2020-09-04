package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	articleRouter := newArticleRouter(db)

	routerGroup := r.Group("/v1/article")
	routerGroup.POST("", articleRouter.post)
	routerGroup.PUT("/:id", articleRouter.updateID)
	routerGroup.GET("/:id", articleRouter.getID)
	routerGroup.GET("", articleRouter.get)
}

func newArticleRouter(db *gorm.DB) articleRouter {
	return articleRouter{
		db: db,
	}
}

type articleRouter struct {
	db *gorm.DB
}

func (ar *articleRouter) post(c *gin.Context) {
	useCase := newUseCaseArticle(c, ar.db)
	result, err := useCase.Create()
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ar *articleRouter) updateID(c *gin.Context) {
	useCase := newUseCaseArticle(c, ar.db)
	err := useCase.UpdateID()
	if err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (ar *articleRouter) getID(c *gin.Context) {
	useCase := newUseCaseArticle(c, ar.db)
	result, err := useCase.GetID()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ar *articleRouter) get(c *gin.Context) {
	useCase := newUseCaseArticle(c, ar.db)
	result, err := useCase.Get()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}
