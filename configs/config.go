package configs

import (
	"os"
	"strings"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Host *configHost

	Gin  *configGin
	Gorm *configGorm

	Env *configEnv

	Validator *configValidator

	reloadFunc []func()
)

func SetReloadFunc(f func()) {
	reloadFunc = append(reloadFunc, f)
}

// Start 開始 Config 設定參數與讀取檔案並轉成 struct
// 預設會抓取執行程式的啟示點資料夾，可用參數調整路徑來源
func Start(sourcePath string) {
	viper.AddConfigPath(getPath(sourcePath))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	Host = newConfigHost()

	Gin = newConfigGin()
	Gorm = newConfigGorm()

	Env = newConfigEnv()

	Validator = newConfigValidator()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		Host.reload()
		Env.reload()

		for _, f := range reloadFunc {
			f()
		}
	})
}

// getPath 預設會抓取執行程式的啟示點資料夾
// e.g. cmd/blog_api 會抓取到 SP_blog
// 可用參數調整路徑來源
func getPath(sourcePath string) string {
	if tools.IsNotEmpty(sourcePath) {
		return sourcePath + "/configs"
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// "/" 切割字串陣列，e.g. 利用陣列 -2 等於往上資料夾兩層
	tempPath := strings.Split(path, "/")
	tempPath = tempPath[:len(tempPath)]

	return strings.Join(tempPath, "/") + "/configs"
}
