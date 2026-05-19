package config

import (
	"errors"
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
			return nil, fmt.Errorf("config not found in %s: %w", path, err)
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

func (c *Config) Validate() error {
	var errs []error

	if c.App.Name == "" {
		errs = append(errs, errors.New("app.name is required"))
	}
	if c.Server.HTTP.Port == 0 {
		errs = append(errs, errors.New("server.http.port is required"))
	}
	if c.Database.Postgres.Host == "" {
		errs = append(errs, errors.New("database.postgres.host is required"))
	}
	if c.Database.Postgres.User == "" {
		errs = append(errs, errors.New("database.postgres.user is required"))
	}
	if c.Database.Postgres.DBName == "" {
		errs = append(errs, errors.New("database.postgres.dbname is required"))
	}

	if len(errs) == 0 {
		return nil
	}

	result := "config validation errors:"
	for _, e := range errs {
		result += "\n  - " + e.Error()
	}
	return errors.New(result)
}

