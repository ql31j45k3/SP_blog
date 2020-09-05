package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
)

const (
	ArticleStatusTag = "articleStatus"

	statusEnable  = 1
	statusDisable = 0
)

func RegisterTagNameFunc(fld reflect.StructField) string {
	return fld.Tag.Get("label")
}

var StatusValidator validator.Func = func(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(int); ok {
		if val == statusDisable || val == statusEnable {
			return true
		}
		return false
	}

	return true
}

var ArticleStatusTranslations validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(ArticleStatusTag, "{0}只有禁用或啟用", true)
}

var ArticleStatusTranslation validator.TranslationFunc = func(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(ArticleStatusTag, fe.Field())
	if err != nil {
		panic(err)
	}
	return t
}
