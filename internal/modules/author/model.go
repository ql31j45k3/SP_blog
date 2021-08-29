package author

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/dig"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

type APIAuthorCond struct {
	dig.In

	R *gin.Engine

	Trans ut.Translator

	DBM *gorm.DB `name:"dbM"`
}

type authorCond struct {
	_ struct{}

	tools.Pagination

	ID uint

	title   string
	content string

	status int
}

func (cond *authorCond) parseArticleID(c *gin.Context) error {
	var err error
	idStr := c.Param("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	return nil
}

func (cond *authorCond) getID(idStr string) (uint, error) {
	if tools.IsEmpty(idStr) {
		return 0, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (cond *authorCond) parseGet(c *gin.Context) error {
	var err error
	idStr := c.Query("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	//nolint:typecheck
	cond.PageIndex, err = tools.Atoi(c.Query("page_index"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	//nolint:typecheck
	cond.PageSize, err = tools.Atoi(c.Query("page_size"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	cond.title = strings.TrimSpace(c.Query("title"))
	cond.content = strings.TrimSpace(c.Query("content"))

	cond.status, err = tools.Atoi(c.Query("status"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	return nil
}

type authors struct {
	_ struct{}

	gorm.Model

	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required,min=10"`

	Status int `json:"status" binding:"authorStatus"`
}

type responseAuthorCreate struct {
	_ struct{}

	ID string `json:"id"`
}

type responseAuthor struct {
	_ struct{}

	tools.Model

	Title   string `json:"title"`
	Content string `json:"content"`

	Status int `json:"status"`
}
