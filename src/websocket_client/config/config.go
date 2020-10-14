package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WebSocketLogPath string `envconfig:"WEBSOCKET_LOG_PATH" default:"/var/log/websocket.log"`
	Host             string `envconfig:"HOST" default:"0.0.0.0:8080"`
	EnvMode          string `envconfig:"ENV_MODE" default:"PRODUCTION"`
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
