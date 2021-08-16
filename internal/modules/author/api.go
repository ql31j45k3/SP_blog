package author

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"gorm.io/gorm"
	"net/http"
)

// RegisterRouter 註冊文章路由器
func RegisterRouter(r *gin.Engine, db *gorm.DB, trans ut.Translator) {
	authorRouter := newAuthorRouter(db, trans)

	routerGroup := r.Group("/v1/author")
	routerGroup.POST("", authorRouter.post)
	routerGroup.PUT("/:id", authorRouter.updateID)
	routerGroup.GET("/:id", authorRouter.getID)
	routerGroup.GET("", authorRouter.get)
}

func newAuthorRouter(db *gorm.DB, trans ut.Translator) authorRouter {
	return authorRouter{
		db:    db,
		trans: trans,
	}
}

type authorRouter struct {
	_ struct{}

	db *gorm.DB

	trans ut.Translator
}

func (ar *authorRouter) post(c *gin.Context) {
	useCase := newUseCaseAuthor(c, ar.db, ar.trans)
	result, err := useCase.Create()
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ar *authorRouter) updateID(c *gin.Context) {
	useCase := newUseCaseAuthor(c, ar.db, ar.trans)
	err := useCase.UpdateID()
	if err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (ar *authorRouter) getID(c *gin.Context) {
	useCase := newUseCaseAuthor(c, ar.db, ar.trans)
	result, err := useCase.GetID()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ar *authorRouter) get(c *gin.Context) {
	useCase := newUseCaseAuthor(c, ar.db, ar.trans)
	result, err := useCase.Get()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, result)
}
