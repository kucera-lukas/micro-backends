package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
	"github.com/kucera-lukas/micro-backends/backend-service/proto/mongo"
	"github.com/kucera-lukas/micro-backends/backend-service/proto/postgres"
)

// NewMessageRepository returns implementation of the controller.Message interface.
func NewMessageRepository(
	mongoClient pbmongo.MessageServiceClient,
	postgresClient pbpostgres.MessageServiceClient,
	rabbitmqClient *rabbitmq.Client,
) controller.Message { //nolint:ireturn
	return &messageRepository{
		mongoClient:    mongoClient,
		postgresClient: postgresClient,
		rabbitmqClient: rabbitmqClient,
	}
}

type messageRepository struct {
	mongoClient    pbmongo.MessageServiceClient
	postgresClient pbpostgres.MessageServiceClient
	rabbitmqClient *rabbitmq.Client
}

func (m messageRepository) Get(
	ctx context.Context,
	id string,
	provider gqlgen.MessageProvider,
) (*gqlgen.Message, error) {
	// TODO implement me
	panic("implement me")
}

func (m messageRepository) List(
	ctx context.Context,
	providers ...gqlgen.MessageProvider,
) ([]*gqlgen.Message, error) {
	var messages []*gqlgen.Message

	for _, provider := range providers {
		if provider == gqlgen.MessageProviderMongo {
			response, err := m.mongoClient.GetMessages(
				ctx,
				&pbmongo.GetMessagesRequest{},
			)
			if err != nil {
				return nil, err
			}

			for _, message := range response.GetMessages() {
				messages = append(messages, &gqlgen.Message{
					ID:       message.GetId(),
					Data:     message.GetData(),
					Created:  message.GetCreated().AsTime(),
					Modified: message.GetModified().AsTime(),
				})
			}
		} else {
			response, err := m.postgresClient.GetMessages(
				ctx,
				&pbpostgres.GetMessagesRequest{},
			)
			if err != nil {
				return nil, err
			}

			for _, message := range response.GetMessages() {
				messages = append(messages, &gqlgen.Message{
					ID:       strconv.Itoa(int(message.GetId())),
					Data:     message.GetData(),
					Created:  message.GetCreated().AsTime(),
					Modified: message.GetModified().AsTime(),
				})
			}
		}
	}

	return messages, nil
}

type Counter interface {
	GetCount() int64
}

func (m messageRepository) Count(
	ctx context.Context,
	providers ...gqlgen.MessageProvider,
) (count int64, err error) {
	var response Counter

	for _, provider := range providers {
		if provider == gqlgen.MessageProviderMongo {
			response, err = m.mongoClient.MessageCount(
				ctx,
				&pbmongo.MessageCountRequest{},
			)
		} else {
			response, err = m.postgresClient.MessageCount(
				ctx,
				&pbpostgres.MessageCountRequest{},
			)
		}

		if err != nil {
			return 0, err
		}

		count += response.GetCount()
	}

	return count, nil
}

func (m messageRepository) Create(
	ctx context.Context,
	data string,
	providers ...gqlgen.MessageProvider,
) (string, error) {
	table := amqp091.Table{"type": rabbitmq.NewMessageKey}

	for _, provider := range providers {
		table[fmt.Sprintf(
			"%s_service",
			strings.ToLower(provider.String()),
		)] = true
	}

	if err := m.rabbitmqClient.Publisher.Publish(
		fmt.Sprintf(`{"data": %q}`, data),
		table,
	); err != nil {
		return "", err
	}

	return "Queued, thanks <3", nil
}

func (m messageRepository) DeliverMessage(
	ctx context.Context,
	delivery amqp091.Delivery,
	messages chan *gqlgen.MessageCreatedPayload,
) {
	var payload *gqlgen.MessageCreatedPayload

	if err := json.Unmarshal(delivery.Body, &payload); err != nil {
		if err := delivery.Nack(false, true); err != nil {
			log.Printf("consume: failed to nack delivery: %v\n", err)
		}
		log.Printf("consume: failed unmarshal delivery body %s: %v\n",
			delivery.Body,
			err,
		)
		return
	}

	messages <- payload

	if err := delivery.Ack(false); err != nil {
		log.Printf("consume: failed to ack delivery: %v\n", err)
	}
}
