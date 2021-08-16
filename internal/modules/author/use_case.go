package author

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseAuthor(c *gin.Context, db *gorm.DB, trans ut.Translator) UseCaseAuthorEr {
	return &author{
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

type author struct {
	_ struct{}

	c  *gin.Context
	db *gorm.DB

	trans ut.Translator
}

func (a *author) Create() (uint, error) {
	var author Author
	if err := tools.BindJSON(a.c, a.trans, &author); err != nil {
		return 0, err
	}

	newRowID, err := a.create(author)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return newRowID, err
	}

	return newRowID, nil
}

func (a *author) UpdateID() error {
	var author Author
	if err := tools.BindJSON(a.c, a.trans, &author); err != nil {
		return err
	}

	ID := a.c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return err
	}

	if err := a.updateID(cond, author); err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return err
	}

	return nil
}

func (a *author) GetID() (ResponseAuthor, error) {
	var responseAuthor ResponseAuthor

	ID := a.c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return responseAuthor, err
	}

	author, err := a.getID(cond)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return responseAuthor, err
	}

	if err := tools.ConvSourceToData(&author, &responseAuthor); err != nil {
		return responseAuthor, err
	}

	return responseAuthor, nil
}

func (a *author) Get() ([]ResponseAuthor, error) {
	var responseAuthors []ResponseAuthor

	cond, err := newAuthorCond(withAuthorPageIndex(a.c.Query("pageIndex")),
		withAuthorPageSize(a.c.Query("pageSize")),
		withAuthorID(a.c.Query("id")),
		withAuthorTitle(a.c.Query("title")),
		withAuthorContent(a.c.Query("content")),
		withAuthorStatus(a.c.Query("status")))
	if err != nil {
		tools.NewReturnError(a.c, http.StatusBadRequest, err)
		return responseAuthors, err
	}

	authors, err := a.get(cond)
	if err != nil {
		tools.IsErrRecordNotFound(a.c, err)
		return responseAuthors, err
	}

	if err := tools.ConvSourceToData(&authors, &responseAuthors); err != nil {
		return responseAuthors, err
	}

	return responseAuthors, nil
}
