package author

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseAuthor(repositoryAuthor repositoryAuthor, db *gorm.DB, trans ut.Translator) useCaseAuthor {
	return &author{
		repositoryAuthor: repositoryAuthor,

		db:    db,
		trans: trans,
	}
}

type useCaseAuthor interface {
	Create(c *gin.Context) (uint, error)
	UpdateID(c *gin.Context) error
	GetID(c *gin.Context) (responseAuthor, error)
	Get(c *gin.Context) ([]responseAuthor, error)
}

type author struct {
	_ struct{}

	repositoryAuthor

	db *gorm.DB

	trans ut.Translator
}

func (a *author) Create(c *gin.Context) (uint, error) {
	var author authors
	if err := tools.BindJSON(c, a.trans, &author); err != nil {
		return 0, err
	}

	newRowID, err := a.repositoryAuthor.Create(a.db, author)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return newRowID, err
	}

	return newRowID, nil
}

func (a *author) UpdateID(c *gin.Context) error {
	var author authors
	if err := tools.BindJSON(c, a.trans, &author); err != nil {
		return err
	}

	ID := c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return err
	}

	if err := a.repositoryAuthor.UpdateID(a.db, cond, author); err != nil {
		tools.IsErrRecordNotFound(c, err)
		return err
	}

	return nil
}

func (a *author) GetID(c *gin.Context) (responseAuthor, error) {
	var responseAuthor responseAuthor

	ID := c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return responseAuthor, err
	}

	author, err := a.repositoryAuthor.GetID(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return responseAuthor, err
	}

	if err := tools.ConvSourceToData(&author, &responseAuthor); err != nil {
		return responseAuthor, err
	}

	return responseAuthor, nil
}

func (a *author) Get(c *gin.Context) ([]responseAuthor, error) {
	var responseAuthors []responseAuthor

	cond, err := newAuthorCond(withAuthorPageIndex(c.Query("page_index")),
		withAuthorPageSize(c.Query("page_size")),
		withAuthorID(c.Query("id")),
		withAuthorTitle(c.Query("title")),
		withAuthorContent(c.Query("content")),
		withAuthorStatus(c.Query("status")))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return responseAuthors, err
	}

	authors, err := a.repositoryAuthor.Get(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return responseAuthors, err
	}

	if err := tools.ConvSourceToData(&authors, &responseAuthors); err != nil {
		return responseAuthors, err
	}

	return responseAuthors, nil
}
