package article

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseArticle(c *gin.Context, db *gorm.DB, trans ut.Translator) UseCaseArticler {
	return &useCaseArticle{
		c:     c,
		db:    db,
		trans: trans,
	}
}

type UseCaseArticler interface {
	Create() (uint, error)
	UpdateID() error
	GetID() (ResponseArticle, error)
	Get() ([]ResponseArticle, error)
	Search() ([]ResponseArticle, error)
}

type useCaseArticle struct {
	c  *gin.Context
	db *gorm.DB

	trans ut.Translator
}

func (uca *useCaseArticle) Create() (uint, error) {
	var article Article
	if err := tools.BindJSON(uca.c, uca.trans, &article); err != nil {
		return 0, err
	}

	newRowID, err := uca.create(article)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return newRowID, err
	}

	return newRowID, nil
}

func (uca *useCaseArticle) UpdateID() error {
	var article Article
	if err := tools.BindJSON(uca.c, uca.trans, &article); err != nil {
		return err
	}

	ID := uca.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return err
	}

	if err := uca.updateID(cond, article); err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return err
	}

	return nil
}

func (uca *useCaseArticle) GetID() (ResponseArticle, error) {
	var responseArticle ResponseArticle

	ID := uca.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return responseArticle, err
	}

	article, err := uca.getID(cond)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return responseArticle, err
	}

	if err := tools.ConvSourceToData(&article, &responseArticle); err != nil {
		return responseArticle, err
	}

	return responseArticle, nil
}

func (uca *useCaseArticle) Get() ([]ResponseArticle, error) {
	var responseArticles []ResponseArticle

	cond, err := newArticleCond(withArticlePageIndex(uca.c.Query("pageIndex")),
		withArticlePageSize(uca.c.Query("pageSize")),
		withArticleID(uca.c.Query("id")),
		withArticleTitle(uca.c.Query("title")),
		withArticleDesc(uca.c.Query("desc")),
		withArticleContent(uca.c.Query("content")),
		withArticleStatus(uca.c.Query("status")))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := uca.get(cond)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return responseArticles, err
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		return responseArticles, err
	}

	return responseArticles, nil
}

func (uca *useCaseArticle) Search() ([]ResponseArticle, error) {
	var responseArticles []ResponseArticle

	cond, err := newSearchCond(withSearchPageIndex(uca.c.Query("pageIndex")),
		withSearchPageSize(uca.c.Query("pageSize")),
		withSearchKeyword(uca.c.Query("keyword")),
		withSearchTags(uca.c.QueryArray("tags")))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := uca.search(cond)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return responseArticles, err
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		return responseArticles, err
	}

	return responseArticles, nil
}
