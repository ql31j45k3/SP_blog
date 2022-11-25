package configs

import (
	"strings"

	"github.com/spf13/viper"
)

func newConfigValidator() *configValidator {
	config := &configValidator{
		locale: viper.GetString("validator.locale"),
	}

	return config
}

type configValidator struct {
	locale string
}

func (c *configValidator) GetLocale() string {
	return strings.ToLower(c.locale)
}
