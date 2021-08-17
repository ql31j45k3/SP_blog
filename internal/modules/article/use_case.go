package article

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseArticle(db *gorm.DB, trans ut.Translator) useCaseArticle {
	return &article{
		db:    db,
		trans: trans,
	}
}

type useCaseArticle interface {
	Create(c *gin.Context) (uint, error)
	UpdateID(c *gin.Context) error
	GetID(c *gin.Context) (responseArticle, error)
	Get(c *gin.Context) ([]responseArticle, error)
	Search(c *gin.Context) ([]responseArticle, error)
}

type article struct {
	_ struct{}

	db *gorm.DB

	trans ut.Translator
}

func (a *article) Create(c *gin.Context) (uint, error) {
	var article articles
	if err := tools.BindJSON(c, a.trans, &article); err != nil {
		return 0, err
	}

	newRowID, err := create(article)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return newRowID, err
	}

	return newRowID, nil
}

func (a *article) UpdateID(c *gin.Context) error {
	var article articles
	if err := tools.BindJSON(c, a.trans, &article); err != nil {
		return err
	}

	ID := c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return err
	}

	if err := a.updateID(cond, article); err != nil {
		tools.IsErrRecordNotFound(c, err)
		return err
	}

	return nil
}

func (a *article) GetID(c *gin.Context) (responseArticle, error) {
	var responseArticle responseArticle

	ID := c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return responseArticle, err
	}

	article, err := a.getID(cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return responseArticle, err
	}

	if err := tools.ConvSourceToData(&article, &responseArticle); err != nil {
		return responseArticle, err
	}

	return responseArticle, nil
}

func (a *article) Get(c *gin.Context) ([]responseArticle, error) {
	var responseArticles []responseArticle

	cond, err := newArticleCond(withArticlePageIndex(c.Query("pageIndex")),
		withArticlePageSize(c.Query("pageSize")),
		withArticleID(c.Query("id")),
		withArticleTitle(c.Query("title")),
		withArticleDesc(c.Query("desc")),
		withArticleContent(c.Query("content")),
		withArticleStatus(c.Query("status")))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := a.get(cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return responseArticles, err
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		return responseArticles, err
	}

	return responseArticles, nil
}

func (a *article) Search(c *gin.Context) ([]responseArticle, error) {
	var responseArticles []responseArticle

	cond, err := newSearchCond(withSearchPageIndex(c.Query("pageIndex")),
		withSearchPageSize(c.Query("pageSize")),
		withSearchKeyword(c.Query("keyword")),
		withSearchTags(c.QueryArray("tags")))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := a.search(cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return responseArticles, err
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		return responseArticles, err
	}

	return responseArticles, nil
}
