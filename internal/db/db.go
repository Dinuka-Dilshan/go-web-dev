package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Address         string
	MaxConns        int32
	MaxConnIdleTime time.Duration
}

func New(ctx context.Context, cfg DBConfig) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.Address)

	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.MaxConnIdleTime = cfg.MaxConnIdleTime
	config.MaxConns = cfg.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	return pool, nil
}
