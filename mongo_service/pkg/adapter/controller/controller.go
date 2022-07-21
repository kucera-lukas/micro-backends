package controller

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/rabbitmq"
)

// Controller holds the Message controller.
type Controller struct {
	Message interface{ Message }
}

// Setup sets up the Controller.
func (c *Controller) Setup(rabbitmqClient *rabbitmq.Client) error {
	ctx := context.Background()

	if err := rabbitmqClient.Consumer.Consume(
		func(delivery amqp091.Delivery) {
			c.Message.NewMessage(ctx, delivery)
		},
		amqp091.Table{
			rabbitmq.ConsumerIdentifier: true,
			"type":                      rabbitmq.NewMessageKey,
			"x-match":                   "all",
		},
	); err != nil {
		return fmt.Errorf("setup: failed to consume messages: %w", err)
	}

	return nil
}
