package repository

import (
	"context"
	"fmt"
	"strconv"

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
	amqpClient *rabbitmq.Client,
) controller.Message { //nolint:ireturn
	return &messageRepository{
		mongoClient:    mongoClient,
		postgresClient: postgresClient,
		amqpClient:     amqpClient,
	}
}

type messageRepository struct {
	mongoClient    pbmongo.MessageServiceClient
	postgresClient pbpostgres.MessageServiceClient
	amqpClient     *rabbitmq.Client
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
	// TODO implement me
	panic("implement me")
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
	if provider == gqlgen.MessageProviderMongo {
		response, err := m.mongoClient.MessageCount(
			ctx,
			&pbmongo.MessageCountRequest{},
		)
		if err != nil {
			return 0, err
		}

		return response.GetCount(), nil
	} else {
		response, err := m.postgresClient.MessageCount(
			ctx,
			&pbpostgres.MessageCountRequest{},
		)
		if err != nil {
			return 0, err
		}

		return response.GetCount(), nil
	}
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
	if provider == gqlgen.MessageProviderMongo {
		response, err := m.mongoClient.NewMessage(
			ctx,
			&pbmongo.NewMessageRequest{Data: data},
		)
		if err != nil {
			return nil, err
		}

		return &gqlgen.Message{
			ID:       response.GetId(),
			Data:     data,
			Provider: provider,
		}, nil
	} else {
		response, err := m.postgresClient.NewMessage(
			ctx,
			&pbpostgres.NewMessageRequest{Data: data},
		)
		if err != nil {
			return nil, err
		}

		return &gqlgen.Message{
			ID:       strconv.Itoa(int(response.GetId())),
			Data:     data,
			Provider: provider,
		}, nil
	}
}

func (m messageRepository) CreateAll(
	ctx context.Context,
	data string,
) (string, error) {
	if err := m.amqpClient.Publish(fmt.Sprintf(`{"data": %q`, data)); err != nil {
		return "", err
	}

	return "Queued, thanks <3", nil
}
