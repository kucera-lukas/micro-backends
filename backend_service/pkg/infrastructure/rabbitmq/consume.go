package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	ConsumerIdentifier = "backend_service"
)

// Consumer for consuming AMQP events.
type Consumer struct {
	conn *amqp091.Connection
}

// NewConsumer returns a new configured Consumer.
func NewConsumer(conn *amqp091.Connection) (*Consumer, error) {
	consumer := &Consumer{conn: conn}

	if err := consumer.setup(); err != nil {
		return nil, err
	}

	return consumer, nil
}

// Consume consumes events from the AMQP exchange.
func (c *Consumer) Consume(
	callback func(delivery amqp091.Delivery),
	tables ...amqp091.Table,
) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("consume: failed to open channel: %w", err)
	}

	queue, err := declareQueue(channel, "")
	if err != nil {
		return fmt.Errorf("consume: failed to declare queue: %w", err)
	}

	for _, table := range tables {
		if err := bindQueue(channel, queue, ExchangeName, table); err != nil {
			return fmt.Errorf(
				"consume: failed to bind queue %q via %T: %w",
				queue.Name,
				table,
				err,
			)
		}
	}

	deliveries, err := channel.Consume(
		queue.Name,
		ConsumerIdentifier,
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
