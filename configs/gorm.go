package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

func newConfigGorm() *configGorm {
	config := &configGorm{
		mode: viper.GetString("gorm.log.mode"),

		username: viper.GetString("database.username"),
		password: viper.GetString("database.password"),
		host:     viper.GetString("database.host"),
		port:     viper.GetString("database.port"),
		dbname:   viper.GetString("database.dbname"),
	}

	config.dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.username, config.password, config.host, config.port, config.dbname)

	return config
}

type configGorm struct {
	_ struct{}

	mode string

	host string
	port string

	username string
	password string

	dbname string

	dsn string
}

func (c *configGorm) GetLogMode() logger.LogLevel {
	if strings.ToLower(c.mode) == "silent" {
		return logger.Silent
	}
	if strings.ToLower(c.mode) == "error" {
		return logger.Error
	}
	if strings.ToLower(c.mode) == "warn" {
		return logger.Warn
	}
	if strings.ToLower(c.mode) == "info" {
		return logger.Info
	}

	return logger.Silent
}

func (c *configGorm) GetDSN() string {
	return c.dsn
}
