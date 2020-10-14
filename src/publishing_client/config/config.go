package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PublishingLogPath string `envconfig:"PUBLISHING_LOG_PATH" default:"publishing.log"`
	Host              string `envconfig:"HOST" default:"0.0.0.0:8080"`
	IntervalTime      int    `envconfig:"INTERVAL_TIME" default:"10"`
	EnvMode           string `envconfig:"ENV_MODE" default:"PRODUCTION"`
}

func GetConfig() (*Config, error) {
	conf := new(Config)
	if err := conf.loadFromEnv(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *Config) loadFromEnv() error {
	return envconfig.Process("", c)
}
