package controller

import (
	"context"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
)

// Message Controller interface.
type Message interface {
	Get(
		ctx context.Context,
		id string,
		provider gqlgen.MessageProvider,
	) (*gqlgen.Message, error)
	List(
		ctx context.Context,
		providers ...gqlgen.MessageProvider,
	) ([]*gqlgen.Message, error)
	Count(
		ctx context.Context,
		providers ...gqlgen.MessageProvider,
	) (int64, error)
	Create(
		ctx context.Context,
		data string,
		providers ...gqlgen.MessageProvider,
	) (string, error)
	DeliverMessage(
		ctx context.Context,
		delivery amqp091.Delivery,
		messages chan *gqlgen.MessageCreatedPayload,
	)
}
