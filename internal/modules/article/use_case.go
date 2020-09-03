package article

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/SP_blog/internal/utils"
	"gorm.io/gorm"
	"net/http"
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

	cond := newArticleCond()
	if err := cond.getID(ID); err != nil {
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

	cond := newArticleCond()
	if err := cond.getID(ID); err != nil {
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
