package article

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/dig"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

type APIArticleCond struct {
	dig.In

	R *gin.Engine

	Trans ut.Translator

	DBM *gorm.DB `name:"dbM"`
}

type articleCondOption func(*articleCond) error

func newArticleCond(opts ...articleCondOption) (*articleCond, error) {
	cond := &articleCond{}

	for _, o := range opts {
		if err := o(cond); err != nil {
			return nil, err
		}
	}

	return cond, nil
}

type articleCond struct {
	_ struct{}

	tools.Pagination

	ID uint

	title   string
	desc    string
	content string

	status int
}

func withArticlePageIndex(pageIndex string) articleCondOption {
	return func(cond *articleCond) error {
		pageIndex, err := tools.Atoi(pageIndex, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.PageIndex = pageIndex
		return nil
	}
}

func withArticlePageSize(pageSize string) articleCondOption {
	return func(cond *articleCond) error {
		pageSize, err := tools.Atoi(pageSize, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.PageSize = pageSize
		return nil
	}
}

func withArticleID(IDStr string) articleCondOption {
	return func(cond *articleCond) error {
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

func withArticleTitle(title string) articleCondOption {
	return func(cond *articleCond) error {
		cond.title = strings.TrimSpace(title)
		return nil
	}
}

func withArticleDesc(desc string) articleCondOption {
	return func(cond *articleCond) error {
		cond.desc = strings.TrimSpace(desc)
		return nil
	}
}

func withArticleContent(content string) articleCondOption {
	return func(cond *articleCond) error {
		cond.content = strings.TrimSpace(content)
		return nil
	}
}

func withArticleStatus(status string) articleCondOption {
	return func(cond *articleCond) error {
		status, err := tools.Atoi(status, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.status = status
		return nil
	}
}

type articles struct {
	_ struct{}

	gorm.Model

	Title   string `json:"title" binding:"required,min=1,max=100"`
	Desc    string `json:"desc" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=10"`

	Status int `json:"status" binding:"articleStatus"`

	ArticleLabel []articleLabels `json:"article_label" gorm:"-"`
}

type articleLabels struct {
	_ struct{}

	gorm.Model

	ArticlesID uint
	Tag        string `binding:"required,min=1,max=100"`
}

type responseArticle struct {
	_ struct{}

	tools.Model

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`

	Status int `json:"status"`
}

type searchCondOption func(*searchCond) error

func newSearchCond(opts ...searchCondOption) (*searchCond, error) {
	cond := &searchCond{}

	for _, o := range opts {
		if err := o(cond); err != nil {
			return nil, err
		}
	}

	return cond, nil
}

type searchCond struct {
	_ struct{}

	tools.Pagination

	keyword string
	tags    []string
}

func withSearchPageIndex(pageIndex string) searchCondOption {
	return func(cond *searchCond) error {
		pageIndex, err := tools.Atoi(pageIndex, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.PageIndex = pageIndex
		return nil
	}
}

func withSearchPageSize(pageSize string) searchCondOption {
	return func(cond *searchCond) error {
		pageSize, err := tools.Atoi(pageSize, tools.DefaultNotAssignInt)
		if err != nil {
			return err
		}

		cond.PageSize = pageSize
		return nil
	}
}

func withSearchKeyword(keyword string) searchCondOption {
	return func(cond *searchCond) error {
		cond.keyword = strings.TrimSpace(keyword)
		return nil
	}
}

func withSearchTags(tags []string) searchCondOption {
	return func(cond *searchCond) error {
		cond.tags = tags
		return nil
	}
}
