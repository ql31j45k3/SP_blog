package article

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseArticle(repositoryArticle repositoryArticle, db *gorm.DB) useCaseArticle {
	return &article{
		repositoryArticle: repositoryArticle,

		db: db,
	}
}

type useCaseArticle interface {
	Create(c *gin.Context, article articles)
	UpdateID(c *gin.Context, id string, article articles)
	GetID(c *gin.Context)
	Get(c *gin.Context)
	Search(c *gin.Context)
}

type article struct {
	_ struct{}

	repositoryArticle

	db *gorm.DB
}

func (a *article) Create(c *gin.Context, article articles) {
	newRowID, err := a.repositoryArticle.Create(a.db, article)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	result := responseArticleCreate{
		ID: fmt.Sprint(newRowID),
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(result))
}

func (a *article) UpdateID(c *gin.Context, id string, article articles) {
	cond, err := newArticleCond(withArticleID(id))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	if err := a.repositoryArticle.UpdateID(a.db, cond, article); err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *article) GetID(c *gin.Context) {
	var responseArticle responseArticle

	ID := c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	article, err := a.repositoryArticle.GetID(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	if err := tools.ConvSourceToData(&article, &responseArticle); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(responseArticle))
}

func (a *article) Get(c *gin.Context) {
	var responseArticles []responseArticle

	cond, err := newArticleCond(withArticlePageIndex(c.Query("page_index")),
		withArticlePageSize(c.Query("page_size")),
		withArticleID(c.Query("id")),
		withArticleTitle(c.Query("title")),
		withArticleDesc(c.Query("desc")),
		withArticleContent(c.Query("content")),
		withArticleStatus(c.Query("status")))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	articles, err := a.repositoryArticle.Get(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(responseArticles))
}

func (a *article) Search(c *gin.Context) {
	var responseArticles []responseArticle

	cond, err := newSearchCond(withSearchPageIndex(c.Query("page_index")),
		withSearchPageSize(c.Query("page_size")),
		withSearchKeyword(c.Query("keyword")),
		withSearchTags(c.QueryArray("tags")))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	articles, err := a.repositoryArticle.Search(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	if err := tools.ConvSourceToData(&articles, &responseArticles); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(responseArticles))
}
