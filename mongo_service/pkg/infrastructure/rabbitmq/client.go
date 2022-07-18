package rabbitmq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName             = "message_events"
	ExchangeKind             = "topic"
	NewMessageRoutingKey     = "new_message"
	CreatedMessageRoutingKey = "created_message"
)

type Client struct {
	Publisher *Publisher
	Consumer  *Consumer
}

// MustNew ensures that a new Client is created and panics if not.
func MustNew(url string) *Client {
	client, err := New(url)
	if err != nil {
		log.Panic(err)
	}

	return client
}

// New tries to create a new Client, returning error if unsuccessful.
func New(url string) (*Client, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf(
			"new: failed to establish connection to rabbitmq: %w",
			err,
		)
	}

	publisher, err := NewPublisher(conn)
	if err != nil {
		return nil, fmt.Errorf("new: failed to create publisher: %w", err)
	}

	consumer, err := NewConsumer(conn)
	if err != nil {
		return nil, fmt.Errorf("new: failed to create consumer: %w", err)
	}

	return &Client{Publisher: publisher, Consumer: consumer}, nil
}

func (c *Client) Close() {
	if err := c.Publisher.conn.Close(); err != nil {
		log.Panicf("close: failed to close rabbitmq connection: %v", err)
	}
}
