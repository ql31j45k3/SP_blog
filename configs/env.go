package configs

import (
	"sync"

	"github.com/spf13/viper"
)

func newConfigEnv() *configEnv {
	viper.SetDefault("system.log.level", "warn")
	viper.SetDefault("system.log.path", "/Users/michael_kao/log/def_log")

	return &configEnv{
		logLevel: viper.GetString("system.log.level"),
		logPath:  viper.GetString("system.log.path"),
	}
}

type configEnv struct {
	_ struct{}

	sync.RWMutex

	logLevel string
	logPath  string
}

func (c *configEnv) reload() {
	c.Lock()
	defer c.Unlock()

	c.logLevel = viper.GetString("system.log.level")
}

func (c *configEnv) GetLogLevel() string {
	c.RLock()
	defer c.RUnlock()

	return c.logLevel
}

func (c *configEnv) GetLogPath() string {
	return c.logPath
}
