package tools

import (
	"errors"
	"reflect"
	"time"
)

// ConvSourceToData sourceData = 來源資料, resultData = 結果資料
// 用反射實作達到動態賦值，不需手動一對一比照欄位給值
// 可支援 []struct or struct，參數需丟入 Ptr 型態
func ConvSourceToData(sourceData, resultData interface{}) error {
	resultDataType := reflect.TypeOf(resultData)
	sourceDataType := reflect.TypeOf(sourceData)

	// 判斷型態兩個要一樣
	if resultDataType.Kind() != sourceDataType.Kind() {
		return errors.New("sourceData and resultData type need same kind")
	}

	// 必須為 Ptr 型態，才可有效修改值
	if resultDataType.Kind() != reflect.Ptr {
		return errors.New("sourceData and resultData need kind is Ptr")
	}

	// 第二次檢查，判斷型態兩個要一樣，因為前一個會是 Ptr 型態
	if resultDataType.Elem().Kind() != sourceDataType.Elem().Kind() {
		return errors.New("sourceData and resultData elem type need same kind")
	}

	// 型態為 Struct 直接進行賦值
	if resultDataType.Elem().Kind() == reflect.Struct {
		convFindFieldAndSetFunc(sourceData, resultData)
		return nil
	}

	// 需判斷是否 Slice 型態
	if resultDataType.Elem().Kind() != reflect.Slice {
		return errors.New("sourceData and resultData need kind is Slice")
	}

	// 用 Elem func 取得 sourceData slice
	sourceDataVale := reflect.ValueOf(sourceData).Elem()
	// 初始化 rspType 型態的 slice
	rspVale := reflect.MakeSlice(resultDataType.Elem(), sourceDataVale.Len(), sourceDataVale.Cap())

	for i := 0; i < sourceDataVale.Len(); i++ {
		// 先取得資料的 Addr 的 Interface 值，才可正常執行 Elem func
		convFindFieldAndSetFunc(sourceDataVale.Index(i).Addr().Interface(), rspVale.Index(i).Addr().Interface())
	}

	// 將 rspVale 賦值成功後的結果，塞回 Client rsp 值
	reflect.ValueOf(resultData).Elem().Set(reflect.ValueOf(rspVale.Interface()))

	return nil
}

func convFindFieldAndSetFunc(sourceData, resultData interface{}) {
	resultDataType := reflect.TypeOf(resultData).Elem()
	resultDataValue := reflect.ValueOf(resultData).Elem()
	sourceDataValue := reflect.ValueOf(sourceData).Elem()

	findFieldAndSet(resultDataType, resultDataValue, sourceDataValue)
}

// findFieldAndSet
// 用遞迴方式處理巢狀 struct 的資料結構
func findFieldAndSet(resultDataType reflect.Type, resultDataValue, sourceDataValue reflect.Value) {
	if resultDataType.Kind() != reflect.Struct {
		if resultDataValue.CanSet() && sourceDataValue.CanSet() {
			reflectSetValue(resultDataType, resultDataValue, sourceDataValue)
		}
		return
	}

	for i := 0; i < resultDataType.NumField(); i++ {
		fieldName := resultDataType.Field(i).Name
		resultDataType2 := resultDataType.Field(i).Type
		resultDataValue2 := resultDataValue.FieldByName(fieldName)
		sourceDataValue2 := sourceDataValue.FieldByName(fieldName)

		if resultDataType.Field(i).Type.Kind() == reflect.Struct {
			rspTypeName := resultDataType.Field(i).Type.String()
			// dig struct 可直接跳過(此struct DI 套件使用)
			if rspTypeName == "dig.In" || rspTypeName == "dig.Out" {
				continue
			}

			findFieldAndSet(resultDataType2, resultDataValue2, sourceDataValue2)
		}

		// 先判斷是否可更改資料，CanSet == false 時異動資料會造成 panic
		if resultDataValue2.CanSet() {
			reflectSetValue(resultDataType2, resultDataValue2, sourceDataValue2)
		}
	}
}

// reflectSetValue 取得 sourceDataValue2 值並賦植給 resultDataValue2
// 目前判斷型態只有 string、Uint、Int、Int64、Float64、time.Time
// TODO: 增加其它類型的 set 實作
func reflectSetValue(resultDataType reflect.Type, resultDataValue2, sourceDataValue2 reflect.Value) {
	kind := resultDataType.Kind()

	switch kind {
	case reflect.String:
		resultDataValue2.SetString(sourceDataValue2.String())
	case reflect.Uint:
		resultDataValue2.SetUint(sourceDataValue2.Uint())
	case reflect.Int, reflect.Int64:
		resultDataValue2.SetInt(sourceDataValue2.Int())
	case reflect.Float64:
		resultDataValue2.SetFloat(sourceDataValue2.Float())
	}

	if resultDataType.String() == "time.Time" {
		resultDataValue2.Set(reflect.ValueOf(sourceDataValue2.Interface().(time.Time)))
	}
}
