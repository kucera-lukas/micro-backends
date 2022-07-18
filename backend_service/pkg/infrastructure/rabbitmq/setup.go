package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func declareQueue(channel *amqp091.Channel) (amqp091.Queue, error) {
	return channel.QueueDeclare(
		"",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
}

func declareExchange(channel *amqp091.Channel, name string, kind string) error {
	return channel.ExchangeDeclare(
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
	topic string,
) error {
	return channel.QueueBind(
		queue.Name,
		topic,
		ExchangeName,
		false, // noWait
		nil,   // args
	)
}

func closeChannel(channel *amqp091.Channel) {
	if err := channel.Close(); err != nil {
		log.Printf("failed to close channel: %v\n", err)
	}
}
