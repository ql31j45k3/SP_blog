package configs

import "github.com/spf13/viper"

func newConfigGin() *configGin {
	config := &configGin{
		mode: viper.GetString("gin.mode"),
	}

	return config
}

type configGin struct {
	mode string
}

func (c *configGin) GetMode() string {
	return c.mode
}
