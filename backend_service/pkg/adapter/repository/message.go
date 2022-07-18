package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
	provider gqlgen.MessageProvider,
) ([]*gqlgen.Message, error) {
	var messages []*gqlgen.Message

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

	return messages, nil
}

func (m messageRepository) ListAll(
	ctx context.Context,
) ([]*gqlgen.Message, error) {
	// TODO implement me
	panic("implement me")
}

func (m messageRepository) Count(
	ctx context.Context,
	provider gqlgen.MessageProvider,
) (int64, error) {
	var count int64

	if provider == gqlgen.MessageProviderMongo {
		response, err := m.mongoClient.MessageCount(
			ctx,
			&pbmongo.MessageCountRequest{},
		)
		if err != nil {
			return 0, err
		}

		count = response.GetCount()
	} else {
		response, err := m.postgresClient.MessageCount(
			ctx,
			&pbpostgres.MessageCountRequest{},
		)
		if err != nil {
			return 0, err
		}

		count = response.GetCount()
	}

	return count, nil
}

func (m messageRepository) CountAll(
	ctx context.Context,
) (int64, error) {
	var result int64

	mongoResponse, err := m.mongoClient.MessageCount(
		ctx,
		&pbmongo.MessageCountRequest{},
	)
	if err != nil {
		return 0, err
	}

	result += mongoResponse.GetCount()

	postgresResponse, err := m.postgresClient.MessageCount(
		ctx,
		&pbpostgres.MessageCountRequest{},
	)
	if err != nil {
		return 0, err
	}

	result += postgresResponse.GetCount()

	return result, nil
}

func (m messageRepository) Create(
	ctx context.Context,
	data string,
	provider gqlgen.MessageProvider,
) (*gqlgen.Message, error) {
	var message *gqlgen.Message

	if provider == gqlgen.MessageProviderMongo {
		response, err := m.mongoClient.NewMessage(
			ctx,
			&pbmongo.NewMessageRequest{Data: data},
		)
		if err != nil {
			return nil, err
		}

		message = &gqlgen.Message{
			ID:       response.GetId(),
			Data:     response.GetData(),
			Created:  response.GetCreated().AsTime(),
			Modified: response.GetModified().AsTime(),
		}
	} else {
		response, err := m.postgresClient.NewMessage(
			ctx,
			&pbpostgres.NewMessageRequest{Data: data},
		)
		if err != nil {
			return nil, err
		}

		message = &gqlgen.Message{
			ID:       strconv.Itoa(int(response.GetId())),
			Data:     response.GetData(),
			Created:  response.GetCreated().AsTime(),
			Modified: response.GetModified().AsTime(),
		}
	}

	if err := m.rabbitmqClient.Publisher.Publish(
		fmt.Sprintf(`
{
    "message": {
        id": %q,
        "data": %q,
        "created": %q,
        "modified": %q
    },
    "provider": %q
}`, message.ID, message.Data, message.Created, message.Modified, provider),
		rabbitmq.CreatedMessageRoutingKey,
	); err != nil {
		return nil, fmt.Errorf(
			"create: failed to publish message: %w",
			err,
		)
	}

	return message, nil
}

func (m messageRepository) CreateAll(
	ctx context.Context,
	data string,
) (string, error) {
	if err := m.rabbitmqClient.Publisher.Publish(
		fmt.Sprintf(`{"data": %q}`, data),
		rabbitmq.NewMessageRoutingKey,
	); err != nil {
		return "", err
	}

	return "Queued, thanks <3", nil
}

func (m messageRepository) Consume(
	ctx context.Context,
	delivery amqp091.Delivery,
	messages chan *gqlgen.MessageCreatedPayload,
) {
	var payload *gqlgen.MessageCreatedPayload

	fmt.Printf("body: %s", delivery.Body)

	if err := json.Unmarshal(delivery.Body, &payload); err != nil {
		if err := delivery.Nack(false, true); err != nil {
			log.Printf("consume: failed to nack delivery: %v\n", err)
			return
		}
	}

	fmt.Printf("payload: %+v, %s", payload.Message, payload.Provider)
	messages <- payload

	if err := delivery.Ack(false); err != nil {
		log.Printf("consume: failed to ack delivery: %v\n", err)
	}
}
