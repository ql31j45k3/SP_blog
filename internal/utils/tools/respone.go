package tools

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"time"
)

// BindJSON 取得 Request body 資料，JSON 資料轉行為 struct
func BindJSON(c *gin.Context, trans ut.Translator, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		var errs []string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err2 := range err.(validator.ValidationErrors) {
				errs = append(errs, err2.Translate(trans))
			}
		} else {
			errs = append(errs, err.Error())
		}

		c.JSON(http.StatusBadRequest, NewResponseError(errs))
		return err
	}

	return nil
}

// NewReturnError 設定回傳錯誤，統一回傳錯誤格式
func NewReturnError(c *gin.Context, code int, err error) {
	messages := []string{err.Error()}
	c.JSON(code, NewResponseError(messages))
}

func NewResponseError(messages []string) ResponseError {
	return ResponseError{
		Message: messages,
	}
}

type ResponseError struct {
	Message []string `json:"messages"`
}

// IsErrRecordNotFound 驗證 SQL 語法執行但查無資料情況，調整 http status
func IsErrRecordNotFound(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NewReturnError(c, http.StatusNotFound, err)
	} else {
		NewReturnError(c, http.StatusInternalServerError, err)
	}
}

// Model 對外回傳基礎欄位
type Model struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

// ConvResponseStruct data = 資料庫資料, rsp = API 回傳資料
// 用反射實作達到動態賦值，不需手動一對一比照欄位給值
// 可支援 []struct or struct，參數需丟入 Ptr 型態
func ConvResponseStruct(data, rsp interface{}) error {
	rspType := reflect.TypeOf(rsp)
	dataType := reflect.TypeOf(data)

	// 判斷型態兩個要一樣
	if rspType.Kind() != dataType.Kind() {
		return errors.New("data and rsp type need same kind")
	}

	// 必須為 Ptr 型態，才可有效修改值
	if rspType.Kind() != reflect.Ptr {
		return errors.New("data and rsp need kind is Ptr")
	}

	// 第二次檢查，判斷型態兩個要一樣，因為前一個會是 Ptr 型態
	if rspType.Elem().Kind() != dataType.Elem().Kind() {
		return errors.New("data and rsp elem type need same kind")
	}

	// 型態為 Struct 直接進行賦值
	if rspType.Elem().Kind() == reflect.Struct {
		convFindFieldAndSetFunc(data, rsp)
		return nil
	}

	// 需判斷是否 Slice 型態
	if rspType.Elem().Kind() != reflect.Slice {
		return errors.New("data and rsp need kind is Slice")
	}

	// 用 Elem func 取得 data slice
	dataVale := reflect.ValueOf(data).Elem()
	// 初始化 rspType 型態的 slice
	rspVale :=  reflect.MakeSlice(rspType.Elem(), dataVale.Len(), dataVale.Cap())

	for i := 0; i < dataVale.Len(); i++ {
		// 先取得資料的 Addr 的 Interface 值，才可正常執行 Elem func
		convFindFieldAndSetFunc(dataVale.Index(i).Addr().Interface(), rspVale.Index(i).Addr().Interface())
	}

	// 將 rspVale 賦值成功後的結果，塞回 Client rsp 值
	reflect.ValueOf(rsp).Elem().Set(reflect.ValueOf(rspVale.Interface()))

	return nil
}

func convFindFieldAndSetFunc(data, rsp interface{}) {
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

		// 先判斷是否可更改資料，CanSet == false 時異動資料會造成 panic
		if rspValue2.CanSet() {
			reflectSetValue(rspType2, rspValue2, dataValue2)
		}
	}
}

// reflectSetValue 取得 dataValue 值並賦植給 rspValue
// 目前判斷型態只有 string、Int、Uint、time.Time
// TODO 增加其它類型的 set 實作
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


