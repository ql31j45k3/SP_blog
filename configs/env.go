package configs

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

func newConfigEnv() *configEnv {
	viper.SetDefault("system.log.level", "warn")
	viper.SetDefault("system.log.path", "/Users/michael_kao/log/def_log")

	viper.SetDefault("system.pprof.status", false)

	viper.SetDefault("system.pprof.block.status", false)
	viper.SetDefault("system.pprof.block.rate", 1000000000)

	viper.SetDefault("system.pprof.mutex.status", false)
	viper.SetDefault("system.pprof.mutex.rate", 1000000000)

	viper.SetDefault("system.shutdown.timeout", time.Duration(10))

	return &configEnv{
		logLevel: viper.GetString("system.log.level"),
		logPath:  viper.GetString("system.log.path"),

		pprofStatus: viper.GetBool("system.pprof.status"),

		pprofBlockStatus: viper.GetBool("system.pprof.block.status"),
		pprofBlockRate:   viper.GetInt("system.pprof.block.rate"),

		pprofMutexStatus: viper.GetBool("system.pprof.mutex.status"),
		pprofMutexRate:   viper.GetInt("system.pprof.mutex.rate"),

		shutdownTimeout: viper.GetDuration("system.shutdown.timeout") * time.Second,
	}
}

type configEnv struct {
	sync.RWMutex

	logLevel string
	logPath  string

	pprofStatus bool

	pprofBlockStatus bool
	pprofBlockRate   int

	pprofMutexStatus bool
	pprofMutexRate   int

	shutdownTimeout time.Duration
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

func (c *configEnv) GetPPROFStatus() bool {
	return c.pprofStatus
}

func (c *configEnv) GetPPROFBlockStatus() bool {
	return c.pprofBlockStatus
}

func (c *configEnv) GetPPROFBlockRate() int {
	return c.pprofBlockRate
}

func (c *configEnv) GetPPROFMutexStatus() bool {
	return c.pprofMutexStatus
}

func (c *configEnv) GetPPROFMutexRate() int {
	return c.pprofMutexRate
}

func (c *configEnv) GetShutdownTimeout() time.Duration {
	return c.shutdownTimeout
}
