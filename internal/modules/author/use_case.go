package author

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
	"net/http"
)

func newUseCaseAuthor(c *gin.Context, db *gorm.DB, trans ut.Translator) UseCaseAuthorEr {
	return &useCaseAuthor{
		c:     c,
		db:    db,
		trans: trans,
	}
}

type UseCaseAuthorEr interface {
	Create() (uint, error)
	UpdateID() error
	GetID() (ResponseAuthor, error)
	Get() ([]ResponseAuthor, error)
}

type useCaseAuthor struct {
	c  *gin.Context
	db *gorm.DB

	trans ut.Translator
}

func (uca *useCaseAuthor) Create() (uint, error) {
	var author Author
	if err := tools.BindJSON(uca.c, uca.trans, &author); err != nil {
		return 0, err
	}

	newRowID, err := uca.create(author)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return newRowID, err
	}

	return newRowID, nil
}

func (uca *useCaseAuthor) UpdateID() error {
	var author Author
	if err := tools.BindJSON(uca.c, uca.trans, &author); err != nil {
		return err
	}

	ID := uca.c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return err
	}

	if err := uca.updateID(cond, author); err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return err
	}

	return nil
}

func (uca *useCaseAuthor) GetID() (ResponseAuthor, error) {
	var responseAuthor ResponseAuthor

	ID := uca.c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return responseAuthor, err
	}

	author, err := uca.getID(cond)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return responseAuthor, err
	}

	if err := tools.ConvSourceToData(&author, &responseAuthor); err != nil {
		return responseAuthor, err
	}

	return responseAuthor, nil
}

func (uca *useCaseAuthor) Get() ([]ResponseAuthor, error) {
	var responseAuthors []ResponseAuthor

	cond, err := newAuthorCond(withAuthorPageIndex(uca.c.Query("pageIndex")),
		withAuthorPageSize(uca.c.Query("pageSize")),
		withAuthorID(uca.c.Query("id")),
		withAuthorTitle(uca.c.Query("title")),
		withAuthorContent(uca.c.Query("content")),
		withAuthorStatus(uca.c.Query("status")))
	if err != nil {
		tools.NewReturnError(uca.c, http.StatusBadRequest, err)
		return responseAuthors, err
	}

	authors, err := uca.get(cond)
	if err != nil {
		tools.IsErrRecordNotFound(uca.c, err)
		return responseAuthors, err
	}

	if err := tools.ConvSourceToData(&authors, &responseAuthors); err != nil {
		return responseAuthors, err
	}

	return responseAuthors, nil
}
