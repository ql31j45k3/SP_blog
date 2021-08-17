package author

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"gorm.io/gorm"
)

// RegisterRouter 註冊文章路由器
func RegisterRouter(r *gin.Engine, db *gorm.DB, trans ut.Translator) {
	author := newUseCaseAuthor(newRepositoryAuthor(), db, trans)
	authorRouter := newAuthorRouter(author)

	routerGroup := r.Group("/v1/author")
	routerGroup.POST("", authorRouter.post)
	routerGroup.PUT("/:id", authorRouter.updateID)
	routerGroup.GET("/:id", authorRouter.getID)
	routerGroup.GET("", authorRouter.get)
}

func newAuthorRouter(author useCaseAuthor) authorRouter {
	return authorRouter{
		author: author,
	}
}

type authorRouter struct {
	_ struct{}

	author useCaseAuthor
}

func (ar *authorRouter) post(c *gin.Context) {
	result, err := ar.author.Create(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ar *authorRouter) updateID(c *gin.Context) {
	err := ar.author.UpdateID(c)
	if err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (ar *authorRouter) getID(c *gin.Context) {
	result, err := ar.author.GetID(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ar *authorRouter) get(c *gin.Context) {
	result, err := ar.author.Get(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}
