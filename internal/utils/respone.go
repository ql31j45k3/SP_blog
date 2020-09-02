package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"time"
)

type Model struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

func ResponeOK(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

func StrconvDataToRsp(data, rsp interface{}) {
	rspType := reflect.TypeOf(rsp).Elem()
	rspValue := reflect.ValueOf(rsp).Elem()
	dataValue := reflect.ValueOf(data).Elem()

	findFieldAndSet(rspType, rspValue, dataValue)
}

// findFieldNameAndSet
// 用遞迴方式處理巢狀 struct 的資料結構
func findFieldAndSet(rspType reflect.Type, rspValue, dataValue reflect.Value) {
	for i := 0; i < rspType.NumField(); i++ {
		fieldName := rspType.Field(i).Name
		rspType2 := rspType.Field(i).Type
		rspValue2 := rspValue.FieldByName(fieldName)
		dataValue2 := dataValue.FieldByName(fieldName)

		if rspType.Field(i).Type.Kind() == reflect.Struct {
			rspTypeName := rspType.Field(i).Type.String()
			// dig struct 可直接跳過(此struct DI 套件使用)
			if rspTypeName == "dig.In" || rspTypeName == "dig.Out" {
				continue
			}

			findFieldAndSet(rspType2, rspValue2, dataValue2)
		}

		reflectSetValue(rspType2, rspValue2, dataValue2)
	}
}

func reflectSetValue(rspType reflect.Type, rspValue, dataValue reflect.Value) {
	kind := rspType.Kind()
	if kind == reflect.String {
		rspValue.SetString(dataValue.Interface().(string))
	}

	if kind == reflect.Int {
		rspValue.SetInt(dataValue.Int())
	}

	if kind == reflect.Uint {
		rspValue.SetUint(dataValue.Uint())
	}

	if rspType.String() == "time.Time" {
		rspValue.Set(reflect.ValueOf(dataValue.Interface().(time.Time)))
	}
}