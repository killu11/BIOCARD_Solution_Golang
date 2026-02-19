package config

import (
	"directory-viewing-service/pkg"

	"github.com/joho/godotenv"
)

const (
	packageName = "config"
)

type Config struct {
	RabbitMQ  DSNConfig
	Postgres  DSNConfig
	Directory *DirectoryCfg
}

func NewConfig() (*Config, error) {
	var config Config
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, pkg.PackageError(packageName, "failed to load env", err)
	}

	if err := parseEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func parseEnv(config *Config) error {
	postgres, err := NewPostgresCfg()
	if err != nil {
		return err
	}

	dir, err := NewDirectoryCfg()
	if err != nil {
		return err
	}

	rm, err := NewRabbitMQConfig()
	if err != nil {
		return err
	}

	config.RabbitMQ = rm
	config.Postgres = postgres
	config.Directory = dir
	return nil
}
