package author

import (
	"bytes"
	"encoding/json"
	"github.com/ql31j45k3/SP_blog/internal/utils/testtools"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	validatorFunc "github.com/ql31j45k3/SP_blog/internal/utils/validator"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	r          *gin.Engine
	db         *gorm.DB
	translator ut.Translator

	debug bool

	authorURL string
)

func start() {
	viper.Set("configFile", "/Users/michael_kao/go/src/github.com/ql31j45k3/SP_blog/configs")

	debug = false
	authorURL = "/v1/author"

	r, db, translator = testtools.Start()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 註冊翻譯器
		err := zhTranslations.RegisterDefaultTranslations(v, translator)
		if err != nil {
			panic(err)
		}

		// 註冊一個函式，獲取struct tag裡自定義的label作為欄位名
		v.RegisterTagNameFunc(validatorFunc.RegisterTagNameFunc)

		// 註冊自定義函式
		if err := v.RegisterValidation(validatorFunc.AuthorStatusTag,
			validatorFunc.AuthorStatusFunc.Validator); err != nil {
			panic(err)
		}

		// 根據提供的標記註冊翻譯
		v.RegisterTranslation(validatorFunc.AuthorStatusTag, translator,
			validatorFunc.AuthorStatusFunc.Translations, validatorFunc.AuthorStatusFunc.Translation)
	}
}

func TestRegisterRouter(t *testing.T) {
	start()

	RegisterRouter(r, db, translator)

	ID := testPost(t)
	testUpdateID(t, ID)
	testGetID(t, ID)

	testGetConditionsID(t, ID)
	testGetConditionsTitle(t)
	testGetConditionsContent(t)
	testGetConditionsStatus(t)
}

func testPost(t *testing.T) string {
	param := make(map[string]interface{})
	param["status"] = 1
	param["title"] = "title unit test post"
	param["content"] = "content unit test post"

	jsonByte, err := json.Marshal(param)
	if err != nil {
		t.Error(err)
		return ""
	}

	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodPost, authorURL, bytes.NewReader(jsonByte))
	if err2 != nil {
		t.Error(err2)
		return ""
	}

	if debug {
		t.Logf("testPost, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusCreated, httpStatus)

	return string(body)
}

func testUpdateID(t *testing.T, ID string) {
	param := make(map[string]interface{})
	param["status"] = 0
	param["title"] = "title unit test update"
	param["content"] = "content unit test update"

	jsonByte, err := json.Marshal(param)
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, _, err2 := testtools.HttptestRequest(r, http.MethodPut, authorURL+"/"+ID, bytes.NewReader(jsonByte))
	if err2 != nil {
		t.Error(err2)
		return
	}

	assert.Equal(t, http.StatusNoContent, httpStatus)
}

func testGetID(t *testing.T, ID string) {
	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodGet, authorURL+"/"+ID, nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetID, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsID(t *testing.T, ID string) {
	urlValues := url.Values{}
	urlValues.Add("id", ID)

	url, err := url.Parse(authorURL + "?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsID, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsTitle(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("title", "title unit test update")

	url, err := url.Parse(authorURL + "?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsTitle, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsContent(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("content", "content unit test update")

	url, err := url.Parse(authorURL + "?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsContent, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsStatus(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("status", "1")

	url, err := url.Parse(authorURL + "?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsStatus, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}
