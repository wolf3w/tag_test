package domain

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServerHost string `envconfig:"SERVER_HOST" required:"true" default:"localhost"`
	ServerPort string `envconfig:"SERVER_PORT" required:"true" default:"8087"`
	RootDir    string `envconfig:"ROOT_DIR" required:"true" default:"./"`
}

func NewFromEnv() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
