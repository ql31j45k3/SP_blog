package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

const (
	ArticleStatusTag = "articleStatus"

	statusEnable  = 1
	statusDisable = 0
)

var (
	locale = ""
)

func SetLocale(l string) {
	locale = l
}

func RegisterTagNameFunc(fld reflect.StructField) string {
	fieldName := strings.ToLower(fld.Name)
	return locale2FieldMap[locale][fieldName]
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
