package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/model"
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

	err := row.Scan(&messageID)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func (r *imageRepository) Count(ctx context.Context) (uint32, error) {
	var count uint32

	row := r.client.QueryRow(ctx, "SELECT count(*) FROM messages;")
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *imageRepository) List(ctx context.Context) ([]model.Message, error) {
	var data []model.Message

	rows, err := r.client.Query(ctx, "SELECT * FROM messages LIMIT 100;")
	if err != nil {
		return nil, fmt.Errorf("count: %w", err)
	}

	for rows.Next() {
		var msg model.Message

		err := rows.Scan(&msg.Id, &msg.Data, &msg.Created, &msg.Modified)
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		data = append(data, msg)
	}

	return data, nil
}
