package configs

import (
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

func newConfigGorm() *configGorm {
	config := &configGorm{
		mode: viper.GetString("gorm.log.mode"),
	}

	return config
}

type configGorm struct {
	_ struct{}

	mode string
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
