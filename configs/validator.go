package configs

import (
	"github.com/spf13/viper"
	"strings"
)

func newConfigValidator() *configValidator {
	config := &configValidator{
		locale: viper.GetString("validator.locale"),
	}

	return config
}

type configValidator struct {
	_ struct{}

	locale string
}

func (c *configValidator) GetLocale() string {
	return strings.ToLower(c.locale)
}
