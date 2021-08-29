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

type authorCondOption func(*authorCond) error

func newAuthorCond(opts ...authorCondOption) (*authorCond, error) {
	cond := &authorCond{}

	for _, o := range opts {
		if err := o(cond); err != nil {
			return nil, err
		}
	}

	return cond, nil
}

type authorCond struct {
	_ struct{}

	tools.Pagination

	ID uint

	title   string
	content string

	status int
}

func withAuthorPageIndex(pageIndex string) authorCondOption {
	return func(cond *authorCond) error {
		pageIndex, err := tools.Atoi(pageIndex, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.PageIndex = pageIndex
		return nil
	}
}

func withAuthorPageSize(pageSize string) authorCondOption {
	return func(cond *authorCond) error {
		pageSize, err := tools.Atoi(pageSize, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.PageSize = pageSize
		return nil
	}
}

func withAuthorID(IDStr string) authorCondOption {
	return func(cond *authorCond) error {
		if tools.IsEmpty(IDStr) {
			return nil
		}

		ID, err := strconv.ParseUint(IDStr, 10, 64)
		if err != nil {
			return err
		}
		cond.ID = uint(ID)

		return nil
	}
}

func withAuthorTitle(title string) authorCondOption {
	return func(cond *authorCond) error {
		cond.title = strings.TrimSpace(title)
		return nil
	}
}

func withAuthorContent(content string) authorCondOption {
	return func(cond *authorCond) error {
		cond.content = strings.TrimSpace(content)
		return nil
	}
}

func withAuthorStatus(status string) authorCondOption {
	return func(cond *authorCond) error {
		status, err := tools.Atoi(status, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.status = status
		return nil
	}
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
