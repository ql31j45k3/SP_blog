package tools

import (
	"reflect"
	"time"
)

type RspError struct {
	Msgs []string `json:"msgs"`
}

func NewRspError(msgs []string) RspError {
	return RspError{
		Msgs: msgs,
	}
}

// Model 對外回傳基礎欄位
type Model struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

// StrconvDataToRsp data = 資料庫資料, rsp = API 回傳資料
// 用反射實作達到動態賦值，不需手動一對一比照欄位給值
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

// reflectSetValue 取得 dataValue 值並賦植給 rspValue
// 目前判斷型態只有 string、Int、Uint、time.Time
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
