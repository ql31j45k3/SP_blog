package configs

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	ConfigHost *configHost
	ConfigDB   *configDB
	ConfigGin  *configGin
	ConfigGorm *configGorm

	ConfigValidator *configValidator
)

// init 控制設定流程
func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.AddConfigPath(path + "/configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	ConfigHost = newConfigHost()
	ConfigDB = newConfigDB()
	ConfigGin = newConfigGin()
	ConfigGorm = newConfigGorm()
	ConfigValidator = newConfigValidator()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		ConfigHost.reload()
	})
}
