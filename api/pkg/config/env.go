package config

import "github.com/kelseyhightower/envconfig"

func LoadFromEnv(cfg interface{}) error {
	return envconfig.Process("", cfg)
}
