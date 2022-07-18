package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
)

func (r *mutationResolver) NewMessage(ctx context.Context, input gqlgen.NewMessageInput) (*gqlgen.NewMessagePayload, error) {
	message, err := r.controller.Message.Create(ctx, input.Data, input.Provider)
	if err != nil {
		return nil, err
	}

	return &gqlgen.NewMessagePayload{
		Message:  message,
		Provider: input.Provider,
	}, nil
}

func (r *mutationResolver) NewGlobalMessage(ctx context.Context, input gqlgen.NewGlobalMessageInput) (*gqlgen.NewGlobalMessagePayload, error) {
	result, err := r.controller.Message.CreateAll(ctx, input.Data)
	if err != nil {
		return nil, err
	}

	return &gqlgen.NewGlobalMessagePayload{Status: result}, nil
}

func (r *queryResolver) Message(ctx context.Context, id string, provider gqlgen.MessageProvider) (*gqlgen.MessagePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Messages(ctx context.Context, provider gqlgen.MessageProvider) (*gqlgen.MessagesPayload, error) {
	messages, err := r.controller.Message.List(ctx, provider)
	if err != nil {
		return nil, err
	}

	return &gqlgen.MessagesPayload{
		Messages: messages,
		Provider: provider,
	}, err
}

func (r *queryResolver) MessageCount(ctx context.Context, provider gqlgen.MessageProvider) (*gqlgen.MessageCountPayload, error) {
	count, err := r.controller.Message.Count(ctx, provider)
	if err != nil {
		return nil, err
	}

	return &gqlgen.MessageCountPayload{
		Count:    int(count),
		Provider: provider,
	}, nil
}

func (r *queryResolver) GlobalMessages(ctx context.Context) (*gqlgen.GlobalMessagesPayload, error) {
	panic("aa")
}

func (r *queryResolver) GlobalMessageCount(ctx context.Context) (*gqlgen.GlobalMessageCountPayload, error) {
	count, err := r.controller.Message.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	return &gqlgen.GlobalMessageCountPayload{Count: int(count)}, nil
}

func (r *subscriptionResolver) MessageCreated(ctx context.Context) (<-chan *gqlgen.MessageCreatedPayload, error) {
	messages := make(chan *gqlgen.MessageCreatedPayload, 1)

	if err := r.rabbitmqClient.Consumer.Consume(
		func(delivery amqp091.Delivery) {
			r.controller.Message.Consume(ctx, delivery, messages)
		},
		rabbitmq.CreatedMessageRoutingKey,
	); err != nil {
		return nil, fmt.Errorf("error consuming messages: %w", err)
	}

	return messages, nil
}
