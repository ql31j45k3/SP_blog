package article

import (
	"github.com/ql31j45k3/SP_blog/internal/utils"
	"gorm.io/gorm"
	"strconv"
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
}

func withArticleID(IDStr string) articleCondOption {
	return func(ac *articleCond) error {
		ID, err := strconv.ParseUint(IDStr, 10, 64)
		if err != nil {
			return err
		}
		ac.ID = uint(ID)

		return nil
	}
}

type Article struct {
	gorm.Model

	Title   string
	Desc    string
	Content string

	Status int
}

type ArticleRsp struct {
	utils.Model

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`

	Status int `json:"status"`
}
