package configs

import (
	"fmt"

	"github.com/spf13/pflag"

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
func Start(sourcePath string) error {
	// 設定自定義 flag to viper
	if err := parseFlag(); err != nil {
		return fmt.Errorf("parseFlag - %w", err)
	}

	viper.AddConfigPath(viper.GetString("configFile"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if isUseVersion() {
			return nil
		}

		return fmt.Errorf("viper.ReadInConfig - %w", err)
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

	return nil
}

func parseFlag() error {
	pflag.Bool("version", false, "version")
	pflag.String("configFile", "", "configFile path")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("viper.BindPFlags - %w", err)
	}

	return nil
}
