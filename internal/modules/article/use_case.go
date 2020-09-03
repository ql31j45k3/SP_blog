package article

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func newUseCaseArticle(c *gin.Context, db *gorm.DB) UseCaseArticler {
	return &useCaseArticle{
		c:  c,
		db: db,
	}
}

type UseCaseArticler interface {
	Create() (uint, error)
	UpdateID() error
	GetID() (ArticleRsp, error)
	Get() ([]ArticleRsp, error)
}

type useCaseArticle struct {
	c  *gin.Context
	db *gorm.DB
}

func (uca *useCaseArticle) Create() (uint, error) {
	var article Article
	uca.c.BindJSON(&article)

	newRowID, err := uca.create(article)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uca.c.String(http.StatusNotFound, err.Error())
		} else {
			uca.c.String(http.StatusInternalServerError, err.Error())
		}
		return newRowID, err
	}

	return newRowID, err
}

func (uca *useCaseArticle) UpdateID() error {
	var article Article
	uca.c.BindJSON(&article)

	ID := uca.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		uca.c.String(http.StatusBadRequest, err.Error())
		return err
	}

	if err := uca.updateID(cond, article); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uca.c.String(http.StatusNotFound, err.Error())
		} else {
			uca.c.String(http.StatusInternalServerError, err.Error())
		}
		return err
	}

	return nil
}

func (uca *useCaseArticle) GetID() (ArticleRsp, error) {
	var articleRsq ArticleRsp

	ID := uca.c.Param("id")

	cond, err := newArticleCond(withArticleID(ID))
	if err != nil {
		uca.c.String(http.StatusBadRequest, err.Error())
		return articleRsq, err
	}

	article, err := uca.getID(cond)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uca.c.String(http.StatusNotFound, err.Error())
		} else {
			uca.c.String(http.StatusInternalServerError, err.Error())
		}
		return articleRsq, err
	}

	utils.StrconvDataToRsp(&article, &articleRsq)
	return articleRsq, nil
}

func (uca *useCaseArticle) Get() ([]ArticleRsp, error) {
	var articleRsqs []ArticleRsp

	var status int
	if uca.c.Query("status") != "" {
		var err error
		status, err = strconv.Atoi(uca.c.Query("status"))
		if err != nil {
			return articleRsqs, err
		}
	}

	cond, err := newArticleCond(withArticleID(uca.c.Query("id")),
		withArticleTitle(uca.c.Query("title")),
		withArticleDesc(uca.c.Query("desc")),
		withArticleContent(uca.c.Query("content")),
		withArticleStatus(status),)
	if err != nil {
		uca.c.String(http.StatusBadRequest, err.Error())
		return articleRsqs, err
	}

	articles, err := uca.get(cond)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uca.c.String(http.StatusNotFound, err.Error())
		} else {
			uca.c.String(http.StatusInternalServerError, err.Error())
		}
		return articleRsqs, err
	}

	articleRsqs = make([]ArticleRsp, len(articles))
	for i := 0; i < len(articles); i++ {
		utils.StrconvDataToRsp(&articles[i], &articleRsqs[i])
	}

	return articleRsqs, nil
}
