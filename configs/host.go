package configs

import (
	"sync"

	"github.com/spf13/viper"
)

func newConfigHost() *configHost {
	config := &configHost{
		spBlogAPIHost: ":" + viper.GetString("api.spBlog.port"),
	}

	return config
}

type configHost struct {
	_ struct{}

	sync.RWMutex

	spBlogAPIHost string
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
