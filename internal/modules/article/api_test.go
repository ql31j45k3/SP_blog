package article

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/ql31j45k3/SP_blog/internal/utils/testtools"
	validatorFunc "github.com/ql31j45k3/SP_blog/internal/utils/validator"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	r          *gin.Engine
	db         *gorm.DB
	translator ut.Translator

	debug bool

	articleURL       string
	searchArticleURL string
)

func start() {
	viper.Set("configFile", "/Users/michael_kao/go/src/github.com/ql31j45k3/SP_blog/configs")

	debug = false
	articleURL = "/v1/article"
	searchArticleURL = "/v1/search/article"

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
		if err := v.RegisterValidation(validatorFunc.ArticleStatusTag,
			validatorFunc.ArticleStatusFunc.Validator); err != nil {
			panic(err)
		}

		// 根據提供的標記註冊翻譯
		v.RegisterTranslation(validatorFunc.ArticleStatusTag, translator,
			validatorFunc.ArticleStatusFunc.Translations, validatorFunc.ArticleStatusFunc.Translation)
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
	testGetConditionsDesc(t)
	testGetConditionsContent(t)
	testGetConditionsStatus(t)

	testSearchArticle(t)
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

	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodPost, articleURL, bytes.NewReader(jsonByte))
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

	httpStatus, _, err2 := testtools.HttptestRequest(r, http.MethodPut, articleURL+"/"+ID, bytes.NewReader(jsonByte))
	if err2 != nil {
		t.Error(err2)
		return
	}

	assert.Equal(t, http.StatusNoContent, httpStatus)
}

func testGetID(t *testing.T, ID string) {
	httpStatus, body, err2 := testtools.HttptestRequest(r, http.MethodGet, articleURL+"/"+ID, nil)
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

	url, err := url.Parse(articleURL + "?" + urlValues.Encode())
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

	url, err := url.Parse(articleURL + "?" + urlValues.Encode())
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

func testGetConditionsDesc(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("desc", "desc unit test update")

	url, err := url.Parse(articleURL + "?" + urlValues.Encode())
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
		t.Logf("testGetConditionsDesc, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}

func testGetConditionsContent(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("content", "content unit test update")

	url, err := url.Parse(articleURL + "?" + urlValues.Encode())
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

	url, err := url.Parse(articleURL + "?" + urlValues.Encode())
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

func testSearchArticle(t *testing.T) {
	urlValues := url.Values{}
	urlValues.Add("keyword", "title")

	url, err := url.Parse(searchArticleURL + "?" + urlValues.Encode())
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
		t.Logf("testSearchArticle, body value = %s", string(body))
	}

	assert.Equal(t, http.StatusOK, httpStatus)
}
