package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/env"
)

// MustNew ensure that a new ent.Client is created and panics if not.
func MustNew(config *env.Config) *pgxpool.Pool {
	client, err := New(context.Background(), config)
	if err != nil {
		log.Panic(err)
	}

	return client
}

func New(ctx context.Context, config *env.Config) (*pgxpool.Pool, error) {
	client, err := pgxpool.Connect(ctx, config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return client, nil
}
