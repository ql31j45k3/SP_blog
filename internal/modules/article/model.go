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

type articleCond struct {
	_ struct{}

	tools.Pagination

	ID uint

	title   string
	desc    string
	content string

	status int
}

func (cond *articleCond) parseArticleID(c *gin.Context) error {
	var err error
	idStr := c.Param("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	return nil
}

func (cond *articleCond) getID(idStr string) (uint, error) {
	if tools.IsEmpty(idStr) {
		return 0, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (cond *articleCond) parseGet(c *gin.Context) error {
	var err error
	idStr := c.Query("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	cond.PageIndex, err = tools.Atoi(c.Query("page_index"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	cond.PageSize, err = tools.Atoi(c.Query("page_size"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	cond.title = strings.TrimSpace(c.Query("title"))
	cond.desc = strings.TrimSpace(c.Query("desc"))
	cond.content = strings.TrimSpace(c.Query("content"))

	cond.status, err = tools.Atoi(c.Query("status"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	return nil
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

type responseArticleCreate struct {
	_ struct{}

	ID string `json:"id"`
}

type responseArticle struct {
	_ struct{}

	tools.Model

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`

	Status int `json:"status"`
}

type searchCond struct {
	_ struct{}

	tools.Pagination

	keyword string
	tags    []string
}

func (cond *searchCond) parseGet(c *gin.Context) error {
	var err error

	cond.PageIndex, err = tools.Atoi(c.Query("page_index"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	cond.PageSize, err = tools.Atoi(c.Query("page_size"), tools.DefaultNotAssignInt)
	if err != nil {
		return err
	}

	cond.keyword = strings.TrimSpace(c.Query("keyword"))

	cond.tags = c.QueryArray("tags")

	return nil
}
