package article

import (
	"net/http"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"

	"github.com/gin-gonic/gin"
)

// RegisterRouter 註冊文章路由器
func RegisterRouter(condAPI APIArticleCond) {
	article := newUseCaseArticle(newRepositoryArticle(), condAPI.DBM, condAPI.Trans)
	articleRouter := newArticleRouter(article)

	routerGroup := condAPI.R.Group("/v1/article")
	routerGroup.POST("", articleRouter.create)
	routerGroup.PUT("/:id", articleRouter.updateID)
	routerGroup.GET("/:id", articleRouter.getID)
	routerGroup.GET("", articleRouter.get)

	condAPI.R.GET("/v1/search/article", articleRouter.search)
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

func (ar *articleRouter) create(c *gin.Context) {
	result, err := ar.article.Create(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(result))
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

	c.JSON(http.StatusOK, tools.NewResponseBasicSuccess(result))
}

func (ar *articleRouter) get(c *gin.Context) {
	result, err := ar.article.Get(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, tools.NewResponseBasicSuccess(result))
}

func (ar *articleRouter) search(c *gin.Context) {
	result, err := ar.article.Search(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, tools.NewResponseBasicSuccess(result))
}
