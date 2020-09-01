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
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()

	for i := 0; i < rspType.NumField(); i++ {
		fieldName := rspType.Field(i).Name

		if rspType.Field(i).Type.Kind() == reflect.Struct {
			rspTypeModel := rspType.Field(i).Type.String()
			if !(rspTypeModel == "utils.Model") {
				continue
			}

			for j := 0; j < dataType.NumField(); j++ {
				dataTypeModel := dataType.Field(j).Type.String()
				if !(dataTypeModel == "gorm.Model") {
					continue
				}

				for n := 0; n < rspType.Field(i).Type.NumField(); n++ {
					fieldName2 := rspType.Field(i).Type.Field(n).Name

					reflectType2 := rspType.Field(i).Type.Field(n).Type
					value2 := rspValue.FieldByName(fieldName).FieldByName(fieldName2)
					dataValueFile2 := dataValue.FieldByName(fieldName).FieldByName(fieldName2)

					reflectSetValue(value2, dataValueFile2, reflectType2)
				}
			}
		}

		reflectType := rspType.Field(i).Type
		value := rspValue.FieldByName(fieldName)
		dataValueFile := dataValue.FieldByName(fieldName)

		reflectSetValue(value, dataValueFile, reflectType)
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
