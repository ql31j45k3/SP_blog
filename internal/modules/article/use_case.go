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
}

type useCaseArticle struct {
	c  *gin.Context
	db *gorm.DB

	trans ut.Translator
}

func (uca *useCaseArticle) Create() (uint, error) {
	var article Article
	if err := uca.bindJSON(&article); err != nil {
		return 0, err
	}

	newRowID, err := uca.create(article)
	if err != nil {
		uca.isErrRecordNotFound(err)
		return newRowID, err
	}

	return newRowID, nil
}

func (uca *useCaseArticle) UpdateID() error {
	var article Article
	if err := uca.bindJSON(&article); err != nil {
		return err
	}

	ID := uca.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		uca.returnError(http.StatusBadRequest, err)
		return err
	}

	if err := uca.updateID(cond, article); err != nil {
		uca.isErrRecordNotFound(err)
		return err
	}

	return nil
}

func (uca *useCaseArticle) GetID() (ResponseArticle, error) {
	var responseArticle ResponseArticle

	ID := uca.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		uca.returnError(http.StatusBadRequest, err)
		return responseArticle, err
	}

	article, err := uca.getID(cond)
	if err != nil {
		uca.isErrRecordNotFound(err)
		return responseArticle, err
	}

	tools.StrconvDataToRsp(&article, &responseArticle)
	return responseArticle, nil
}

func (uca *useCaseArticle) Get() ([]ResponseArticle, error) {
	var responseArticles []ResponseArticle

	cond, err := newArticleCond(withArticleID(uca.c.Query("id")),
		withArticleTitle(uca.c.Query("title")),
		withArticleDesc(uca.c.Query("desc")),
		withArticleContent(uca.c.Query("content")),
		withArticleStatus(uca.c.Query("status")))
	if err != nil {
		uca.returnError(http.StatusBadRequest, err)
		return responseArticles, err
	}

	articles, err := uca.get(cond)
	if err != nil {
		uca.isErrRecordNotFound(err)
		return responseArticles, err
	}

	responseArticles = make([]ResponseArticle, len(articles))
	for i := 0; i < len(articles); i++ {
		tools.StrconvDataToRsp(&articles[i], &responseArticles[i])
	}

	return responseArticles, nil
}
