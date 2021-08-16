package configs

import (
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
	"strings"
)

func newConfigGorm() *configGorm {
	config := &configGorm{
		mode: viper.GetString("gorm.logmode"),
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
