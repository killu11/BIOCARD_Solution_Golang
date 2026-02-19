package config

import (
	"directory-viewing-service/pkg"
	"fmt"

	"github.com/caarlos0/env/v10"
)

type DSNConfig interface {
	DSN() string
}

type PostgresCfg struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD"`
	DB       string `env:"DB_NAME" envDefault:"processing_files_db"`
	Port     uint   `env:"DB_PORT" envDefault:"5432"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"none"`
}

func NewPostgresCfg() (*PostgresCfg, error) {
	var postgres PostgresCfg
	if err := env.Parse(&postgres); err != nil {
		return nil, pkg.PackageError(packageName, "failed to parse postgres settings", err)
	}
	return &postgres, nil
}

func (c PostgresCfg) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
		c.SSLMode,
	)
}
