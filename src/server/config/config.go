package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddress string `envconfig:"LISTEN_ADDRESS" default:"0.0.0.0:8080"`
	ServerLogPath string `envconfig:"SERVER_LOG_PATH" default:"/var/log/server.log"`
	TemporaryFile string `envconfig:"TEMPORARY_FILE" default:"data.messages.tmp"`
	EnvMode       string `envconfig:"ENV_MODE" default:"DEVELOPMENT"`
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
