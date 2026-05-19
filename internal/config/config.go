package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const configPath = "configs"

func Load(cfgPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	path := cfgPath
	if path == "" {
		path = configPath
	}
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found in %s: %w", path, err)
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
