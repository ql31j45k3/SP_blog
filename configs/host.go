package configs

import (
	"sync"

	"github.com/spf13/viper"
)

func newConfigHost() *configHost {
	config := &configHost{
		spBlogApisHost: ":" + viper.GetString("host.spBlogApisHost"),
	}

	return config
}

type configHost struct {
	_ struct{}

	sync.RWMutex

	spBlogApisHost string
}

func (c *configHost) reload() {
	c.Lock()
	defer c.Unlock()

	c.spBlogApisHost = ":" + viper.GetString("host.spBlogApisHost")
}

func (c *configHost) GetSPBlogApisHost() string {
	c.RLock()
	defer c.RUnlock()

	return c.spBlogApisHost
}
