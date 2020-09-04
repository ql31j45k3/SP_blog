package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func newConfigDB() *configDB {
	config := &configDB{
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

type configDB struct {
	host string
	port string

	username string
	password string

	dbname string

	dsn string
}

func (c *configDB) GetDSN() string {
	return c.dsn
}
