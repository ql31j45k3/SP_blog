package article

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"

	"github.com/gin-gonic/gin"
)

// RegisterRouter 註冊文章路由器
func RegisterRouter(condAPI APIArticleCond) {
	article := newUseCaseArticle(newRepositoryArticle(), condAPI.DBM)
	articleRouter := articleRouter{
		article: article,
		trans:   condAPI.Trans,
	}

	routerGroup := condAPI.R.Group("/v1/article")
	routerGroup.POST("", articleRouter.create)
	routerGroup.PUT("/:id", articleRouter.updateID)
	routerGroup.GET("/:id", articleRouter.getID)
	routerGroup.GET("", articleRouter.get)

	condAPI.R.GET("/v1/search/article", articleRouter.search)
}

type articleRouter struct {
	_ struct{}

	article useCaseArticle

	trans ut.Translator
}

func (ar *articleRouter) create(c *gin.Context) {
	var article articles
	if err := tools.BindJSON(c, ar.trans, &article); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	result, err := ar.article.Create(c, article)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(result))
}

func (ar *articleRouter) updateID(c *gin.Context) {
	id := c.Param("id")

	var article articles
	if err := tools.BindJSON(c, ar.trans, &article); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	err := ar.article.UpdateID(c, id, article)
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
