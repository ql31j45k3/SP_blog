package testtools

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"

	utilsDriver "github.com/ql31j45k3/SP_blog/internal/utils/driver"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/ql31j45k3/SP_blog/configs"
	validatorFunc "github.com/ql31j45k3/SP_blog/internal/utils/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Start() (*gin.Engine, *gorm.DB, ut.Translator, error) {
	if err := configs.Start(getPath()); err != nil {
		return nil, nil, nil, err
	}

	validatorFunc.Start()

	r := gin.Default()

	db, err := utilsDriver.NewMysql(configs.Gorm.GetMasterHost(), configs.Gorm.GetMasterUsername(), configs.Gorm.GetMasterPassword(),
		configs.Gorm.GetMasterDBName(), configs.Gorm.GetMasterPort(), configs.Gorm.GetLogMode(),
		configs.Gorm.GetMasterMaxIdle(), configs.Gorm.GetMasterMaxOpen(), configs.Gorm.GetMasterMaxLifetime())
	if err != nil {
		panic(err)
	}

	locale := configs.Validator.GetLocale()
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator(locale)
	// 設定語言地區
	validatorFunc.SetLocale(locale)

	return r, db, trans, nil
}

func getPath() string {
	path, err2 := os.Getwd()
	if err2 != nil {
		panic(err2)
	}

	return path[0:strings.Index(path, "SP_blog")] + "SP_blog"
}

// HttptestRequest 根據特定請求 URL 和參數 param
func HttptestRequest(r *gin.Engine, method, uri string, reader io.Reader) (int, []byte, error) {
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
