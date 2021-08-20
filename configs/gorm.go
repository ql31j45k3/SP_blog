package configs

import (
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
		dbName:   viper.GetString("database.dbName"),
	}

	return config
}

type configGorm struct {
	_ struct{}

	mode string

	host string
	port string

	username string
	password string

	dbName string
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

func (c *configGorm) GetHost() string {
	return c.host
}

func (c *configGorm) GetPort() string {
	return c.port
}

func (c *configGorm) GetUsername() string {
	return c.username
}

func (c *configGorm) GetPassword() string {
	return c.password
}

func (c *configGorm) GetDBName() string {
	return c.dbName
}
