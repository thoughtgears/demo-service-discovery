package cfg

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port         string `envconfig:"PORT" default:"8080"`
	DiscoveryURL string `envconfig:"DISCOVERY_URL"`
	BackendURL   string `envconfig:"BACKEND_URL"`
	Service      string `envconfig:"K_SERVICE"`
	Environment  string `envconfig:"ENVIRONMENT" default:"dev"`
}

func NewConfig() (*Config, error) {
	var config Config
	envconfig.MustProcess("", &config)

	if config.DiscoveryURL == "" && config.BackendURL == "" {
		return nil, errors.New("DISCOVERY_URL or BACKEND_URL must be set")
	}

	return &config, nil
}
