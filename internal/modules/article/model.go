package article

import (
	"github.com/ql31j45k3/SP_blog/internal/utils"
	"gorm.io/gorm"
	"strconv"
)

func newArticleCond() articleCond {
	return articleCond{}
}

type articleCond struct {
	ID uint
}

func (ac *articleCond) getID(IDStr string) error {
	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		return err
	}
	ac.ID = uint(ID)

	return nil
}

type Article struct {
	gorm.Model

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`

	State int `json:"state"`
}

type ArticleRsp struct {
	utils.Model

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`

	State int `json:"state"`
}
