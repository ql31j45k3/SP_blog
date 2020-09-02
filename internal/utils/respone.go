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

		reflectSetValue(rspValue2, dataValue2, rspType2)
	}
}

func reflectSetValue(value, dataValueFile reflect.Value, reflectType reflect.Type) {
	kind := reflectType.Kind()
	if kind == reflect.String {
		value.SetString(dataValueFile.Interface().(string))
	}

	if kind == reflect.Int {
		value.SetInt(dataValueFile.Int())
	}

	if kind == reflect.Uint {
		value.SetUint(dataValueFile.Uint())
	}

	if reflectType.String() == "time.Time" {
		value.Set(reflect.ValueOf(dataValueFile.Interface().(time.Time)))
	}
}
