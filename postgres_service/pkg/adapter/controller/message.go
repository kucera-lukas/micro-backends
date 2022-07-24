package controller

import (
	"context"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/model"
)

// Message Controller interface.
type Message interface {
	Get(ctx context.Context, messageID string) (*model.Message, error)
	Create(ctx context.Context, data string) (*model.Message, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]*model.Message, error)
	NewMessage(ctx context.Context, delivery amqp091.Delivery)
}
