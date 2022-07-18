package controller

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/rabbitmq"
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
			c.Message.Consume(ctx, delivery)
		},
		rabbitmq.NewMessageRoutingKey,
	); err != nil {
		return fmt.Errorf("setup: failed to consume messages: %w", err)
	}

	return nil
}
