package article

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/ql31j45k3/SP_blog/configs"
	validatorFunc "github.com/ql31j45k3/SP_blog/internal/utils/validator"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	r          *gin.Engine
	db         *gorm.DB
	translator ut.Translator

	debug bool
)

func start() {
	debug = false

	path, err2 := os.Getwd()
	if err2 != nil {
		panic(err2)
	}
	path = path[0:strings.Index(path, "SP_blog")] + "SP_blog"

	configs.Start(path)
	validatorFunc.Start()

	r = gin.Default()

	var err error
	db, err = gorm.Open(mysql.Open(configs.ConfigDB.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(configs.ConfigGorm.GetLogMode()),
	})
	if err != nil {
		panic(err)
	}

	locale := configs.ConfigValidator.GetLocale()
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator(locale)
	// 設定語言地區
	validatorFunc.SetLocale(locale)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 註冊翻譯器
		err := zhTranslations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			panic(err)
		}

		// 註冊一個函式，獲取struct tag裡自定義的label作為欄位名
		v.RegisterTagNameFunc(validatorFunc.RegisterTagNameFunc)

		// 註冊自定義函式
		if err := v.RegisterValidation(validatorFunc.ArticleStatusTag,
			validatorFunc.ArticleStatusFunc.Validator); err != nil {
			panic(err)
		}

		// 根據提供的標記註冊翻譯
		v.RegisterTranslation(validatorFunc.ArticleStatusTag, trans,
			validatorFunc.ArticleStatusFunc.Translations, validatorFunc.ArticleStatusFunc.Translation)
	}
	translator = trans
}

// httptestRequest 根據特定請求 URL 和參數 param
func httptestRequest(r *gin.Engine, method, uri string, reader io.Reader) (int, []byte, error) {
	req := httptest.NewRequest(method, uri, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return 0, nil, err
	}

	return w.Code, body, nil
}

func TestRegisterRouter(t *testing.T) {
	start()

	RegisterRouter(r, db, translator)

	ID := testPost(t)
	testUpdateID(t, ID)
	testGetID(t, ID)

	testGetConditionsID(t, ID)
	testGetConditionsTitle(t)
	testGetConditionsDesc(t)
	testGetConditionsContent(t)
	testGetConditionsStatus(t)
}

func testPost(t *testing.T) string {
	param := make(map[string]interface{})
	param["status"] = 1
	param["title"] = "title unit test post"
	param["desc"] = "desc unit test post"
	param["content"] = "content unit test post"

	jsonByte, err := json.Marshal(param)
	if err != nil {
		t.Error(err)
		return ""
	}

	httpStatus, body, err2 := httptestRequest(r, http.MethodPost, "/v1/article", bytes.NewReader(jsonByte))
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
	param["desc"] = "desc unit test update"
	param["content"] = "content unit test update"

	jsonByte, err := json.Marshal(param)
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, _, err2 := httptestRequest(r, http.MethodPut, "/v1/article/"+ID, bytes.NewReader(jsonByte))
	if err2 != nil {
		t.Error(err2)
		return
	}

	assert.Equal(t, http.StatusNoContent, httpStatus)
}

func testGetID(t *testing.T, ID string) {
	httpStatus, body, err2 := httptestRequest(r, http.MethodGet, "/v1/article/"+ID, nil)
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

	url, err := url.Parse("/v1/article?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := httptestRequest(r, http.MethodGet, url.String(), nil)
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

	url, err := url.Parse("/v1/article?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := httptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsTitle, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsDesc(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("desc", "desc unit test update")

	url, err := url.Parse("/v1/article?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := httptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsDesc, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsContent(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("content", "content unit test update")

	url, err := url.Parse("/v1/article?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := httptestRequest(r, http.MethodGet, url.String(), nil)
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

	url, err := url.Parse("/v1/article?" + urlValues.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	httpStatus, body, err2 := httptestRequest(r, http.MethodGet, url.String(), nil)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if debug {
		t.Logf("testGetConditionsStatus, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}
