package validator

import "github.com/ql31j45k3/SP_blog/internal/utils/validator/zh"

var (
	locale2FieldMap map[string]map[string]string
)

func init()  {
	locale2FieldMap = make(map[string]map[string]string)

	locale2FieldMap["zh"] = zh.NewField2Name()
}
