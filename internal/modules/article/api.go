package article

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	articleRouter := newArticleRouter(db)

	routerGroup := r.Group("/v1/article")
	routerGroup.POST("", articleRouter.post)
	routerGroup.GET("/:id", articleRouter.getID)
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
	result, err := useCase.Post()
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ar *articleRouter) getID(c *gin.Context) {
	useCase := newUseCaseArticle(c, ar.db)
	result, err := useCase.GetID()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}
