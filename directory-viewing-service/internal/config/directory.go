package config

import (
	"directory-viewing-service/pkg"
	"errors"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
)

var ErrInvalidPath = errors.New("invalid i/o paths, check .env")

type DirectoryCfg struct {
	In            string        `env:"WORK_DIR"`
	Out           string        `env:"OUT_DIR"`
	WatchInterval time.Duration `env:"WATCH_INTERVAL" envDefault:"5m"`
}

func NewDirectoryCfg() (*DirectoryCfg, error) {
	var cfg DirectoryCfg
	if err := env.Parse(&cfg); err != nil {
		return nil, pkg.PackageError(packageName, "parse env to dirCfg", err)
	}
	if cfg.In == "" || cfg.Out == "" {
		return nil, ErrInvalidPath
	}

	if strings.HasSuffix(cfg.In, "/") || strings.HasSuffix(cfg.Out, "/") {
		return nil, ErrInvalidPath
	}
	return &cfg, nil
}
