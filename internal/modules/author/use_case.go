package author

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newUseCaseAuthor(repositoryAuthor repositoryAuthor, db *gorm.DB) useCaseAuthor {
	return &author{
		repositoryAuthor: repositoryAuthor,

		db: db,
	}
}

type useCaseAuthor interface {
	Create(c *gin.Context, author authors)
	UpdateID(c *gin.Context, id string, author authors)
	GetID(c *gin.Context)
	Get(c *gin.Context)
}

type author struct {
	_ struct{}

	repositoryAuthor

	db *gorm.DB

	trans ut.Translator
}

func (a *author) Create(c *gin.Context, author authors) {
	newRowID, err := a.repositoryAuthor.Create(a.db, author)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	result := responseAuthorCreate{
		ID: fmt.Sprint(newRowID),
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(result))
}

func (a *author) UpdateID(c *gin.Context, id string, author authors) {
	cond, err := newAuthorCond(withAuthorID(id))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	if err := a.repositoryAuthor.UpdateID(a.db, cond, author); err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *author) GetID(c *gin.Context) {
	var responseAuthor responseAuthor

	ID := c.Param("id")

	cond, err := newAuthorCond(withAuthorID(ID))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	author, err := a.repositoryAuthor.GetID(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	if err := tools.ConvSourceToData(&author, &responseAuthor); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(responseAuthor))
}

func (a *author) Get(c *gin.Context) {
	var responseAuthors []responseAuthor

	cond, err := newAuthorCond(withAuthorPageIndex(c.Query("page_index")),
		withAuthorPageSize(c.Query("page_size")),
		withAuthorID(c.Query("id")),
		withAuthorTitle(c.Query("title")),
		withAuthorContent(c.Query("content")),
		withAuthorStatus(c.Query("status")))
	if err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	authors, err := a.repositoryAuthor.Get(a.db, cond)
	if err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	if err := tools.ConvSourceToData(&authors, &responseAuthors); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, tools.NewResponseBasicSuccess(responseAuthors))
}
