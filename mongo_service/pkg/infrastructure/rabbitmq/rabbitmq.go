package rabbitmq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/env"
)

const (
	QueueName = "messages"
)

type Client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

// MustNew ensures that a new Client is created and panics if not.
func MustNew(config *env.Config) *Client {
	client, err := New(config)
	if err != nil {
		log.Panic(err)
	}

	return client
}

// New tries to create a new Client, returning error if unsuccessful.
func New(config *env.Config) (*Client, error) {
	conn, err := amqp091.Dial(config.RabbitMQURI)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to establish connection to rabbitmq: %w",
			err,
		)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		QueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (c *Client) Publish(body string) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
	}

	if err = channel.Publish(
		"",           // exchange
		c.queue.Name, // key
		true,         // mandatory
		false,        // immediate
		msg,
	); err != nil {
		return err
	}

	return nil
}

func (c *Client) Consume(callback func(delivery amqp091.Delivery)) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	deliveries, err := channel.Consume(
		c.queue.Name,
		"",
		true,  // autoAck
		true,  // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			callback(delivery)
		}
	}()

	return nil
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		log.Panicf("failed to close rabbitmq connection: %v", err)
	}

	if err := c.channel.Close(); err != nil {
		log.Panicf("failed to close rabbitmq channel: %v", err)
	}
}
