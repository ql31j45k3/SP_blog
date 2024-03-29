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
	article useCaseArticle

	trans ut.Translator
}

func (ar *articleRouter) create(c *gin.Context) {
	var article articles
	if err := tools.BindJSON(c, ar.trans, &article); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.article.Create(c, article)
}

func (ar *articleRouter) updateID(c *gin.Context) {
	cond := articleCond{}
	if err := cond.parseArticleID(c); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	var article articles
	if err := tools.BindJSON(c, ar.trans, &article); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.article.UpdateID(c, cond, article)
}

func (ar *articleRouter) getID(c *gin.Context) {
	cond := articleCond{}
	if err := cond.parseArticleID(c); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.article.GetID(c, cond)
}

func (ar *articleRouter) get(c *gin.Context) {
	cond := articleCond{}
	if err := cond.parseGet(c); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.article.Get(c, cond)
}

func (ar *articleRouter) search(c *gin.Context) {
	cond := searchCond{}
	if err := cond.parseGet(c); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.article.Search(c, cond)
}
