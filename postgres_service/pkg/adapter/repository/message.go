package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/controller"
)

// NewMessageRepository returns implementation of the controller.Message interface.
func NewMessageRepository(client *pgxpool.Pool) controller.Message { //nolint:ireturn
	return &imageRepository{client: client}
}

type imageRepository struct {
	client *pgxpool.Pool
}

func (r *imageRepository) Create(
	ctx context.Context,
	data string,
) (uint32, error) {
	var messageID uint32

	row := r.client.QueryRow(ctx, "INSERT INTO messages (data) VALUES ($1) RETURNING messages.id;", data)

	err := row.Scan(messageID)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func (r *imageRepository) Count(ctx context.Context) (uint32, error) {
	var count uint32

	rows, err := r.client.Query(ctx, "SELECT count(*) FROM messages;")
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	err = rows.Scan(count)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *imageRepository) List(ctx context.Context) (string, error) {
	var data string

	rows, err := r.client.Query(ctx, "WITH cte AS (SELECT * FROM messages) SELECT to_json(cte) FROM cte;")
	if err != nil {
		return "", fmt.Errorf("count: %w", err)
	}

	err = rows.Scan(data)
	if err != nil {
		return "", fmt.Errorf("count: %w", err)
	}

	return data, nil
}
