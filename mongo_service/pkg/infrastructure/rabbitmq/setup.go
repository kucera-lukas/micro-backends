package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func declareQueue(
	channel *amqp091.Channel,
	name string,
) (amqp091.Queue, error) {
	return channel.QueueDeclare( // nolint:wrapcheck
		name,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
}

func declareExchange(channel *amqp091.Channel, name, kind string) error {
	return channel.ExchangeDeclare( // nolint:wrapcheck
		name,
		kind,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)
}

func bindQueue(
	channel *amqp091.Channel,
	queue amqp091.Queue,
	exchange string,
	args amqp091.Table,
) error {
	return channel.QueueBind( // nolint:wrapcheck
		queue.Name,
		"", // 'headers' exchange ignores the routing key
		exchange,
		false, // noWait
		args,
	)
}

func closeChannel(channel *amqp091.Channel) {
	if err := channel.Close(); err != nil {
		log.Printf("failed to close channel: %v\n", err)
	}
}
