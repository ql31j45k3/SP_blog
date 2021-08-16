package configs

import "github.com/spf13/viper"

func newConfigHost() *configHost {
	config := &configHost{
		spBlogApisHost: ":" + viper.GetString("host.sp_blog_apis_host"),
	}

	return config
}

type configHost struct {
	_ struct{}

	spBlogApisHost string
}

func (c *configHost) reload() {
	c.spBlogApisHost = ":" + viper.GetString("host.sp_blog_apis_host")
}

func (c *configHost) GetSPBlogApisHost() string {
	return c.spBlogApisHost
}
