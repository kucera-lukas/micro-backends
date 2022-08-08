package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
	pbmongo "github.com/kucera-lukas/micro-backends/backend-service/proto/mongo"
	pbpostgres "github.com/kucera-lukas/micro-backends/backend-service/proto/postgres"
)

// NewMessageRepository returns implementation of the controller.Message interface.
func NewMessageRepository( // nolint:ireturn
	mongoClient pbmongo.MessageServiceClient,
	postgresClient pbpostgres.MessageServiceClient,
	rabbitmqClient *rabbitmq.Client,
) controller.Message {
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

type messageGetter interface {
	GetId() string
	GetData() string
	GetCreated() *timestamppb.Timestamp
	GetModified() *timestamppb.Timestamp
}

func (m messageRepository) Get(
	ctx context.Context,
	messageID string,
	provider gqlgen.MessageProvider,
) (message *gqlgen.Message, err error) {
	var response messageGetter

	if provider == gqlgen.MessageProviderMongo {
		response, err = m.mongoClient.GetMessage(
			ctx,
			&pbmongo.GetMessageRequest{Id: messageID},
		)
	} else {
		response, err = m.postgresClient.GetMessage(
			ctx,
			&pbpostgres.GetMessageRequest{Id: messageID},
		)
	}

	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &gqlgen.Message{
		ID:       response.GetId(),
		Data:     response.GetData(),
		Created:  response.GetCreated().AsTime(),
		Modified: response.GetModified().AsTime(),
	}, nil
}

func (m messageRepository) List(
	ctx context.Context,
	sortField gqlgen.MessageSortField,
	reverse bool,
	providers ...gqlgen.MessageProvider,
) (messages []*gqlgen.Message, err error) {
	for _, provider := range providers {
		if provider == gqlgen.MessageProviderMongo {
			response, err := m.mongoClient.GetMessages(
				ctx,
				&pbmongo.GetMessagesRequest{},
			)
			if err != nil {
				return nil, fmt.Errorf("list: %w", err)
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
				return nil, fmt.Errorf("list: %w", err)
			}

			for _, message := range response.GetMessages() {
				messages = append(messages, &gqlgen.Message{
					ID:       message.GetId(),
					Data:     message.GetData(),
					Created:  message.GetCreated().AsTime(),
					Modified: message.GetModified().AsTime(),
				})
			}
		}
	}

	sort.Slice(messages, getMessageSortFunc(messages, sortField, reverse))

	return messages, nil
}

type counter interface {
	GetCount() int64
}

func (m messageRepository) Count(
	ctx context.Context,
	providers ...gqlgen.MessageProvider,
) (count int64, err error) {
	var response counter

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
			return 0, fmt.Errorf("count: %w", err)
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
	if len(providers) == 0 {
		return "You didn't provide any providers :(", nil
	}

	table := amqp091.Table{"type": rabbitmq.NewMessageKey}

	for _, provider := range providers {
		table[fmt.Sprintf(
			"%s_service",
			strings.ToLower(provider.String()),
		)] = true
	}

	if err := m.rabbitmqClient.Publisher.Publish(
		ctx,
		fmt.Sprintf(`{"data": %q}`, data),
		table,
	); err != nil {
		return "", fmt.Errorf("create: %w", err)
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

		log.Printf("consume: failed to unmarshal delivery body %s: %v\n",
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

func getMessageSortFunc( // nolint:cyclop
	messages []*gqlgen.Message,
	field gqlgen.MessageSortField,
	reverse bool,
) func(i, j int) bool {
	if reverse {
		switch field {
		case gqlgen.MessageSortFieldID:
			return func(i, j int) bool {
				return messages[i].ID > messages[j].ID
			}
		case gqlgen.MessageSortFieldData:
			return func(i, j int) bool {
				return messages[i].Data > messages[j].Data
			}
		case gqlgen.MessageSortFieldCreated:
			return func(i, j int) bool {
				return messages[i].Created.After(messages[j].Created)
			}
		case gqlgen.MessageSortFieldModified:
			return func(i, j int) bool {
				return messages[i].Modified.After(messages[j].Modified)
			}
		default:
			return func(i, j int) bool {
				return messages[i].ID > messages[j].ID
			}
		}
	} else {
		switch field {
		case gqlgen.MessageSortFieldID:
			return func(i, j int) bool {
				return messages[i].ID < messages[j].ID
			}
		case gqlgen.MessageSortFieldData:
			return func(i, j int) bool {
				return messages[i].Data < messages[j].Data
			}
		case gqlgen.MessageSortFieldCreated:
			return func(i, j int) bool {
				return messages[i].Created.Before(messages[j].Created)
			}
		case gqlgen.MessageSortFieldModified:
			return func(i, j int) bool {
				return messages[i].Modified.Before(messages[j].Modified)
			}
		default:
			return func(i, j int) bool {
				return messages[i].ID < messages[j].ID
			}
		}
	}
}
