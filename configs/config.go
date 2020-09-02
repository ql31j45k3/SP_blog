package configs

import (
	"os"

	"github.com/spf13/viper"
)

var (
	ConfigHost *configHost
	ConfigDB *configDB
)

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
}
