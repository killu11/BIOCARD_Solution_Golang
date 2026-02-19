package persistence

import "C"
import (
	"context"
	"directory-viewing-service/internal/config"
	"directory-viewing-service/pkg"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const packageName = "persistence"

func NewPostgresConnections(c config.DSNConfig) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(c.DSN())
	if err != nil {
		return nil, pkg.PackageError(packageName, "failed to open db connection", err)
	}
	cfg.MaxConns = 50                      // Максимальное количество соединений
	cfg.MinConns = 10                      // Минимальное количество соединений
	cfg.MaxConnLifetime = 30 * time.Minute // Максимальное время жизни соединения
	cfg.MaxConnIdleTime = 30 * time.Minute // Максимальное время простоя

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, pkg.PackageError(packageName, "ping conns", err)
	}
	return pool, nil
}
