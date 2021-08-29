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
	UpdateID(c *gin.Context, cond authorCond, author authors)
	GetID(c *gin.Context, cond authorCond)
	Get(c *gin.Context, cond authorCond)
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

func (a *author) UpdateID(c *gin.Context, cond authorCond, author authors) {
	if err := a.repositoryAuthor.UpdateID(a.db, cond, author); err != nil {
		tools.IsErrRecordNotFound(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *author) GetID(c *gin.Context, cond authorCond) {
	var responseAuthor responseAuthor

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

func (a *author) Get(c *gin.Context, cond authorCond) {
	var responseAuthors []responseAuthor

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
