package article

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRouter 註冊文章路由器
func RegisterRouter(r *gin.Engine, db *gorm.DB, trans ut.Translator) {
	article := newUseCaseArticle(newRepositoryArticle(), db, trans)
	articleRouter := newArticleRouter(article)

	routerGroup := r.Group("/v1/article")
	routerGroup.POST("", articleRouter.post)
	routerGroup.PUT("/:id", articleRouter.updateID)
	routerGroup.GET("/:id", articleRouter.getID)
	routerGroup.GET("", articleRouter.get)

	r.GET("/v1/search/article", articleRouter.search)
}

func newArticleRouter(article useCaseArticle) articleRouter {
	return articleRouter{
		article: article,
	}
}

type articleRouter struct {
	_ struct{}

	article useCaseArticle
}

func (ar *articleRouter) post(c *gin.Context) {
	result, err := ar.article.Create(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ar *articleRouter) updateID(c *gin.Context) {
	err := ar.article.UpdateID(c)
	if err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (ar *articleRouter) getID(c *gin.Context) {
	result, err := ar.article.GetID(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ar *articleRouter) get(c *gin.Context) {
	result, err := ar.article.Get(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ar *articleRouter) search(c *gin.Context) {
	result, err := ar.article.Search(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}
