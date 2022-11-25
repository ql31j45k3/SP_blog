package configs

import (
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

func newConfigGorm() *configGorm {
	viper.SetDefault("database.mysql.master.conn.maxIdle", 10)
	viper.SetDefault("database.mysql.master.conn.maxOpen", 100)
	viper.SetDefault("database.mysql.master.conn.maxLifetime", time.Duration(600))

	config := &configGorm{
		mode: viper.GetString("gorm.log.mode"),

		username: viper.GetString("database.mysql.master.username"),
		password: viper.GetString("database.mysql.master.password"),
		host:     viper.GetString("database.mysql.master.host"),
		port:     viper.GetString("database.mysql.master.port"),
		dbName:   viper.GetString("database.mysql.master.dbName"),

		maxIdle:     viper.GetInt("database.mysql.master.conn.maxIdle"),
		maxOpen:     viper.GetInt("database.mysql.master.conn.maxOpen"),
		maxLifetime: viper.GetDuration("database.mysql.master.conn.maxLifetime") * time.Second,
	}

	return config
}

type configGorm struct {
	mode string

	host string
	port string

	username string
	password string

	dbName string

	maxIdle     int
	maxOpen     int
	maxLifetime time.Duration
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

func (c *configGorm) GetMasterHost() string {
	return c.host
}

func (c *configGorm) GetMasterPort() string {
	return c.port
}

func (c *configGorm) GetMasterUsername() string {
	return c.username
}

func (c *configGorm) GetMasterPassword() string {
	return c.password
}

func (c *configGorm) GetMasterDBName() string {
	return c.dbName
}

func (c *configGorm) GetMasterMaxIdle() int {
	return c.maxIdle
}

func (c *configGorm) GetMasterMaxOpen() int {
	return c.maxOpen
}

func (c *configGorm) GetMasterMaxLifetime() time.Duration {
	return c.maxLifetime
}
