package configs

import (
	"sync"

	"github.com/spf13/viper"
)

func newConfigHost() *configHost {
	config := &configHost{
		spBlogAPIHost: ":" + viper.GetString("api.spBlog.port"),
		pprofAPIHost:  ":" + viper.GetString("api.pprof.port"),
	}

	return config
}

type configHost struct {
	sync.RWMutex

	spBlogAPIHost string
	pprofAPIHost  string
}

func (c *configHost) reload() {
	c.Lock()
	defer c.Unlock()

	c.spBlogAPIHost = ":" + viper.GetString("api.spBlog.port")
}

func (c *configHost) GetSPBlogAPIHost() string {
	c.RLock()
	defer c.RUnlock()

	return c.spBlogAPIHost
}

func (c *configHost) GetPPROFAPIHost() string {
	return c.pprofAPIHost
}
