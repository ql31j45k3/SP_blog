package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/ql31j45k3/SP_blog/internal/utils/validator/zh"
	"reflect"
	"strings"
)

const (
	ArticleStatusTag = "articleStatus"
	AuthorStatusTag  = "authorStatus"

	statusEnable  = 1
	statusDisable = 0
)

var (
	locale2FieldMap map[string]map[string]string

	locale = ""

	ArticleStatusFunc articleStatusFunc
	AuthorStatusFunc  authorStatusFunc
)

func Start() {
	locale2FieldMap = make(map[string]map[string]string)

	locale2FieldMap["zh"] = zh.NewField2Name()

	ArticleStatusFunc = articleStatusFunc{}
	AuthorStatusFunc = authorStatusFunc{}
}

// SetLocale 設定語言地區
func SetLocale(l string) {
	locale = l
}

// RegisterTagNameFunc 註冊欄位對應轉譯的文字
func RegisterTagNameFunc(fld reflect.StructField) string {
	fieldName := strings.ToLower(fld.Name)
	return locale2FieldMap[locale][fieldName]
}

type articleStatusFunc struct {
	_ struct{}
}

// Validator 提供驗證 ArticleStatus 資料正確性 func
func (asf *articleStatusFunc) Validator(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(int); ok {
		if val == statusDisable || val == statusEnable {
			return true
		}
		return false
	}

	return true
}

// Translationsfunc 提供 ArticleStatus 錯誤訊息格式
func (asf *articleStatusFunc) Translations(ut ut.Translator) error {
	return ut.Add(ArticleStatusTag, "{0}只有禁用或啟用", true)
}

// Translation 提供 ArticleStatus 翻譯功能
func (asf *articleStatusFunc) Translation(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(ArticleStatusTag, fe.Field())
	if err != nil {
		panic(err)
	}
	return t
}

type authorStatusFunc struct {
	_ struct{}
}

// Validator 提供驗證 AuthorStatus 資料正確性 func
func (asf *authorStatusFunc) Validator(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(int); ok {
		if val == statusDisable || val == statusEnable {
			return true
		}
		return false
	}

	return true
}

// Translationsfunc 提供 AuthorStatus 錯誤訊息格式
func (asf *authorStatusFunc) Translations(ut ut.Translator) error {
	return ut.Add(AuthorStatusTag, "{0}只有禁用或啟用", true)
}

// Translation 提供 AuthorStatus 翻譯功能
func (asf *authorStatusFunc) Translation(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(AuthorStatusTag, fe.Field())
	if err != nil {
		panic(err)
	}
	return t
}
