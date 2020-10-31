package article

import (
	"strconv"
	"strings"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

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
	ID uint

	title   string
	desc    string
	content string

	status int
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

type Article struct {
	gorm.Model

	Title   string `binding:"required,min=1,max=100"`
	Desc    string `binding:"required,min=1,max=255"`
	Content string `binding:"required,min=10"`

	Status int `binding:"articleStatus"`

	ArticleLabel []ArticleLabel `gorm:"-"`
}

type ArticleLabel struct {
	gorm.Model

	ArticlesID uint
	Tag string `binding:"required,min=1,max=100"`
}

type ResponseArticle struct {
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
	keyword string
	tags []string
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
