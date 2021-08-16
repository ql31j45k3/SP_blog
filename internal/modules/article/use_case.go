package article

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseArticle(c *gin.Context, db *gorm.DB, trans ut.Translator) useCaseArticle {
	return &article{
		c:     c,
		db:    db,
		trans: trans,
	}
}

type useCaseArticle interface {
	Create() (uint, error)
	UpdateID() error
	GetID() (responseArticle, error)
	Get() ([]responseArticle, error)
	Search() ([]responseArticle, error)
}

type article struct {
	_ struct{}

	c  *gin.Context
	db *gorm.DB

	trans ut.Translator
}

func (a *article) Create() (uint, error) {
	var article articles
	if err := tools.BindJSON(a.c, a.trans, &article); err != nil {
		return 0, err
	}

	newRowID, err := a.create(article)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return newRowID, err
	}

	return newRowID, nil
}

func (a *article) UpdateID() error {
	var article articles
	if err := tools.BindJSON(a.c, a.trans, &article); err != nil {
		return err
	}

	ID := a.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return err
	}

	if err := a.updateID(cond, article); err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return err
	}

	return nil
}

func (a *article) GetID() (responseArticle, error) {
	var responseArticle responseArticle

	ID := a.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return responseArticle, err
	}

	article, err := a.getID(cond)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return responseArticle, err
	}

	if err := tools.ConvSourceToData(&article, &responseArticle); err != nil {
		return responseArticle, err
	}

	return responseArticle, nil
}

func (a *article) Get() ([]responseArticle, error) {
	var responseArticles []responseArticle

	cond, err := newArticleCond(withArticlePageIndex(a.c.Query("pageIndex")),
		withArticlePageSize(a.c.Query("pageSize")),
		withArticleID(a.c.Query("id")),
		withArticleTitle(a.c.Query("title")),
		withArticleDesc(a.c.Query("desc")),
		withArticleContent(a.c.Query("content")),
		withArticleStatus(a.c.Query("status")))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := a.get(cond)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return responseArticles, err
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		return responseArticles, err
	}

	return responseArticles, nil
}

func (a *article) Search() ([]responseArticle, error) {
	var responseArticles []responseArticle

	cond, err := newSearchCond(withSearchPageIndex(a.c.Query("pageIndex")),
		withSearchPageSize(a.c.Query("pageSize")),
		withSearchKeyword(a.c.Query("keyword")),
		withSearchTags(a.c.QueryArray("tags")))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := a.search(cond)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return responseArticles, err
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		return responseArticles, err
	}

	return responseArticles, nil
}
