package author

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"

	"github.com/gin-gonic/gin"
)

// RegisterRouter 註冊文章路由器
func RegisterRouter(condAPI APIAuthorCond) {
	author := newUseCaseAuthor(newRepositoryAuthor(), condAPI.DBM)
	authorRouter := authorRouter{
		author: author,

		trans: condAPI.Trans,
	}

	routerGroup := condAPI.R.Group("/v1/author")
	routerGroup.POST("", authorRouter.create)
	routerGroup.PUT("/:id", authorRouter.updateID)
	routerGroup.GET("/:id", authorRouter.getID)
	routerGroup.GET("", authorRouter.get)
}

type authorRouter struct {
	_ struct{}

	author useCaseAuthor

	trans ut.Translator
}

func (ar *authorRouter) create(c *gin.Context) {
	var author authors
	if err := tools.BindJSON(c, ar.trans, &author); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.author.Create(c, author)
}

func (ar *authorRouter) updateID(c *gin.Context) {
	id := c.Param("id")

	var author authors
	if err := tools.BindJSON(c, ar.trans, &author); err != nil {
		tools.NewReturnError(c, http.StatusBadRequest, err)
		return
	}

	ar.author.UpdateID(c, id, author)
}

func (ar *authorRouter) getID(c *gin.Context) {
	ar.author.GetID(c)
}

func (ar *authorRouter) get(c *gin.Context) {
	ar.author.Get(c)
}
