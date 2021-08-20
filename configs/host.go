package configs

import (
	"sync"

	"github.com/spf13/viper"
)

func newConfigHost() *configHost {
	config := &configHost{
		spBlogApisHost: ":" + viper.GetString("host.sp_blog_apis_host"),
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

	c.spBlogApisHost = ":" + viper.GetString("host.sp_blog_apis_host")
}

func (c *configHost) GetSPBlogApisHost() string {
	c.RLock()
	defer c.RUnlock()

	return c.spBlogApisHost
}
