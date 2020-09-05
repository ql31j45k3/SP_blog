package article

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http"

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
	GetID() (ArticleRsp, error)
	Get() ([]ArticleRsp, error)
}

type useCaseArticle struct {
	c  *gin.Context
	db *gorm.DB

	trans ut.Translator
}

func (uca *useCaseArticle) Create() (uint, error) {
	var article Article
	if err := uca.c.ShouldBindJSON(&article); err != nil {
		var errs []string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err2 := range err.(validator.ValidationErrors) {
				errs = append(errs, err2.Translate(uca.trans))
			}
		} else {
			errs = append(errs, err.Error())
		}

		uca.c.JSON(http.StatusBadRequest,
			tools.RspError{
				Msgs: errs,
			})
		return 0, err
	}

	newRowID, err := uca.create(article)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uca.c.String(http.StatusNotFound, err.Error())
		} else {
			uca.c.String(http.StatusInternalServerError, err.Error())
		}
		return newRowID, err
	}

	return newRowID, nil
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

	tools.StrconvDataToRsp(&article, &articleRsq)
	return articleRsq, nil
}

func (uca *useCaseArticle) Get() ([]ArticleRsp, error) {
	var articleRsqs []ArticleRsp

	status, err := tools.Atoi(uca.c.Query("status"), tools.DefaultNotAssignInt)
	if err != nil {
		uca.c.String(http.StatusBadRequest, err.Error())
		return articleRsqs, err
	}

	cond, err := newArticleCond(withArticleID(uca.c.Query("id")),
		withArticleTitle(uca.c.Query("title")),
		withArticleDesc(uca.c.Query("desc")),
		withArticleContent(uca.c.Query("content")),
		withArticleStatus(status))
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
		tools.StrconvDataToRsp(&articles[i], &articleRsqs[i])
	}

	return articleRsqs, nil
}
