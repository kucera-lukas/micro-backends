package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

// Publisher for publishing AMQP events.
type Publisher struct {
	conn *amqp091.Connection
}

// NewPublisher returns a new configured Publisher.
func NewPublisher(conn *amqp091.Connection) (*Publisher, error) {
	publisher := &Publisher{conn: conn}

	if err := publisher.setup(); err != nil {
		return nil, err
	}

	return publisher, nil
}

// Publish a message to the AMQP exchange.
func (p *Publisher) Publish(body string, headers amqp091.Table) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("publish: failed to open channel: %w", err)
	}

	defer closeChannel(channel)

	if err := channel.Publish(
		ExchangeName,
		"",    // 'headers' exchange ignores the routing key
		false, // mandatory
		false, // immediate
		amqp091.Publishing{ // nolint:exhaustivestruct
			ContentType: "application/json",
			Body:        []byte(body),
			Headers:     headers,
		},
	); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

func (p *Publisher) setup() error {
	channel, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("setup: failed to open channel: %w", err)
	}

	defer closeChannel(channel)

	return declareExchange(channel, ExchangeName, ExchangeKind)
}
