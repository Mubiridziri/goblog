package config

import (
	"errors"
	"fmt"
	"goblog/internal/utils"
)

const DatabaseDsn = "DATABASE_DSN"

type ConfigLoader struct {
}

type Config struct {
	Database Database
}

type Database struct {
	DSN string
}

func (loader *ConfigLoader) createConfig() *Config {
	return &Config{
		Database: Database{
			DSN: utils.GetEnvStr(DatabaseDsn, ""),
		},
	}
}

func (loader *ConfigLoader) LoadConfig() (*Config, error) {
	var cfg *Config

	cfg = loader.createConfig()
	if err := loader.validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (loader *ConfigLoader) validate(cfg *Config) error {
	if cfg.Database.DSN == "" {
		return loader.createNotNullEnvError(DatabaseDsn)
	}

	return nil
}

func (loader *ConfigLoader) createNotNullEnvError(envName string) error {
	return errors.New(fmt.Sprintf("env variable %v cannot be null", envName))
}
