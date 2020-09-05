package article

import (
	"strconv"
	"strings"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

type articleCondOption func(*articleCond) error

func newArticleCond(opts ...articleCondOption) (*articleCond, error) {
	ac := &articleCond{}

	for _, o := range opts {
		if err := o(ac); err != nil {
			return ac, err
		}
	}

	return ac, nil
}

type articleCond struct {
	ID uint

	title   string
	desc    string
	content string

	status int
}

func withArticleID(IDStr string) articleCondOption {
	return func(ac *articleCond) error {
		if tools.IsEmpty(IDStr) {
			return nil
		}

		ID, err := strconv.ParseUint(IDStr, 10, 64)
		if err != nil {
			return err
		}
		ac.ID = uint(ID)

		return nil
	}
}

func withArticleTitle(title string) articleCondOption {
	return func(ac *articleCond) error {
		ac.title = strings.TrimSpace(title)
		return nil
	}
}

func withArticleDesc(desc string) articleCondOption {
	return func(ac *articleCond) error {
		ac.desc = strings.TrimSpace(desc)
		return nil
	}
}

func withArticleContent(content string) articleCondOption {
	return func(ac *articleCond) error {
		ac.content = strings.TrimSpace(content)
		return nil
	}
}

func withArticleStatus(status int) articleCondOption {
	return func(ac *articleCond) error {
		ac.status = status
		return nil
	}
}

type Article struct {
	gorm.Model

	Title   string `binding:"required,min=1,max=100" label:"標題"`
	Desc    string `binding:"required,min=1,max=255" label:"描敘"`
	Content string `binding:"required,min=10" label:"內容"`

	Status int `binding:"articleStatus" label:"狀態"`
}

type ArticleRsp struct {
	tools.Model

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`

	Status int `json:"status"`
}
