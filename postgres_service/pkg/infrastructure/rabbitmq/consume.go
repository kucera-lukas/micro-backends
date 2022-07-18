package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	consumerIdentifier = "backend_service"
)

// Consumer for consuming AMQP events
type Consumer struct {
	conn *amqp091.Connection
}

// NewConsumer returns a new configured Consumer.
func NewConsumer(conn *amqp091.Connection) (*Consumer, error) {
	consumer := &Consumer{conn: conn}

	err := consumer.setup()
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// Consume consumes events from the AMQP exchange.
func (c *Consumer) Consume(
	callback func(delivery amqp091.Delivery),
	topics ...string,
) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("consume: failed to open channel: %w", err)
	}

	queue, err := declareQueue(channel)
	if err != nil {
		return fmt.Errorf("consume: failed to declare queue: %w", err)
	}

	for _, topic := range topics {
		if err := bindQueue(channel, queue, topic); err != nil {
			return fmt.Errorf(
				"consume: failed to bind queue %q to the %q topic: %w",
				queue.Name,
				topic,
				err,
			)
		}
	}

	deliveries, err := channel.Consume(
		queue.Name,
		consumerIdentifier,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"consume: failed to start consuming deliveries: %w",
			err,
		)
	}

	go func() {
		for delivery := range deliveries {
			callback(delivery)
		}
	}()

	return nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("setup: failed to open channel: %w", err)
	}

	defer closeChannel(channel)

	return declareExchange(channel, ExchangeName, ExchangeKind)
}
