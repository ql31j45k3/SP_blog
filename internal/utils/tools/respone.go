package tools

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
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

// Pagination 查詢分頁欄位
type Pagination struct {
	PageIndex int
	PageSize  int
}

func (p *Pagination) GetOffset() int {
	if p.PageIndex == DefaultNotAssignInt {
		return 0
	}

	return (p.PageIndex - 1) * p.PageSize
}

func (p *Pagination) GetRowCount() int {
	if p.PageSize == DefaultNotAssignInt || p.PageSize == 0 {
		return 25
	}

	return p.PageSize
}
