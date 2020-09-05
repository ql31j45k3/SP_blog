package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"net/http"
)

func (uca *useCaseArticle) bindJSON(article *Article) error {
	if err := uca.c.BindJSON(article); err != nil {
		var errs []string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err2 := range err.(validator.ValidationErrors) {
				errs = append(errs, err2.Translate(uca.trans))
			}
		} else {
			errs = append(errs, err.Error())
		}

		uca.c.JSON(http.StatusBadRequest,
			tools.RspError{
				Msgs: errs,
			})
		return err
	}

	return nil
}
