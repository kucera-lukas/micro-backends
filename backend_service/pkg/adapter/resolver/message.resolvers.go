package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
)

func (r *mutationResolver) NewMessage(ctx context.Context, input gqlgen.NewMessageInput) (*gqlgen.NewMessagePayload, error) {
	status, err := r.controller.Message.Create(
		ctx,
		input.Data,
		input.Providers...,
	)
	if err != nil {
		return nil, err
	}

	return &gqlgen.NewMessagePayload{
		Status:    status,
		Providers: input.Providers,
	}, nil
}

func (r *queryResolver) Message(ctx context.Context, id string, provider gqlgen.MessageProvider) (*gqlgen.MessagePayload, error) {
	message, err := r.controller.Message.Get(ctx, id, provider)
	if err != nil {
		return nil, err
	}

	return &gqlgen.MessagePayload{
		Message:  message,
		Provider: provider,
	}, nil
}

func (r *queryResolver) Messages(ctx context.Context, providers []gqlgen.MessageProvider, sortField gqlgen.MessageSortField, reverse bool) (*gqlgen.MessagesPayload, error) {
	messages, err := r.controller.Message.List(
		ctx,
		sortField,
		reverse,
		providers...,
	)
	if err != nil {
		return nil, err
	}

	return &gqlgen.MessagesPayload{
		Messages:  messages,
		Providers: providers,
	}, err
}

func (r *queryResolver) MessageCount(ctx context.Context, providers []gqlgen.MessageProvider) (*gqlgen.MessageCountPayload, error) {
	count, err := r.controller.Message.Count(ctx, providers...)
	if err != nil {
		return nil, err
	}

	return &gqlgen.MessageCountPayload{
		Count:     int(count),
		Providers: providers,
	}, nil
}

func (r *subscriptionResolver) MessageCreated(ctx context.Context) (<-chan *gqlgen.MessageCreatedPayload, error) {
	messages := make(chan *gqlgen.MessageCreatedPayload, 1)

	if err := r.rabbitmqClient.Consumer.Consume(
		func(delivery amqp091.Delivery) {
			r.controller.Message.DeliverMessage(ctx, delivery, messages)
		},
		amqp091.Table{
			"provider": strings.ToLower(gqlgen.MessageProviderMongo.String()),
			"type":     rabbitmq.CreatedMessageKey,
			"x-match":  "all",
		}, amqp091.Table{
			"provider": strings.ToLower(gqlgen.MessageProviderPostgres.String()),
			"type":     rabbitmq.CreatedMessageKey,
			"x-match":  "all",
		},
	); err != nil {
		return nil, fmt.Errorf(
			"message_created: error consuming messages: %w",
			err,
		)
	}

	return messages, nil
}
