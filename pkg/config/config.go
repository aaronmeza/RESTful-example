package config

import (
	"github.com/kelseyhightower/envconfig"
)

type MyRetailConfig struct {
	ApiAddress string `envconfig:"PRODUCT_URL"`
	ApiQuery   string `envconfig:"PRODUCT_QUERY"`
}

func Parse() (*MyRetailConfig, error) {

	c := &MyRetailConfig{}
	err := envconfig.Process("", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
