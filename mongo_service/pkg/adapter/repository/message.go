package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/rabbitmq"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/model"
)

const (
	databaseName   = "micro_backends"
	collectionName = "messages"
	providerName   = "MONGO"
	listLimit      = 100
)

type messageRepository struct {
	collection     *mongo.Collection
	rabbitmqClient *rabbitmq.Client
}

// NewMessageRepository returns implementation of the
// controller.Message interface.
func NewMessageRepository( // nolint:ireturn
	mongoClient *mongo.Client,
	rabbitmqClient *rabbitmq.Client,
) controller.Message {
	return &messageRepository{
		collection: mongoClient.
			Database(databaseName).
			Collection(collectionName),
		rabbitmqClient: rabbitmqClient,
	}
}

func (r *messageRepository) Get(
	ctx context.Context,
	messageID string,
) (*model.Message, error) {
	var message model.Message

	objectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, fmt.Errorf(
			"get: faild to parse id %q: %w",
			messageID,
			err,
		)
	}

	if err := r.collection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: objectID}},
	).Decode(&message); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &message, nil
}

func (r *messageRepository) Create(
	ctx context.Context,
	data string,
) (*model.Message, error) {
	now := time.Now()

	message := &model.Message{
		ID:       primitive.NewObjectID(),
		Data:     data,
		Created:  now,
		Modified: now,
	}

	if _, err := r.collection.InsertOne(ctx, message); err != nil {
		return message, fmt.Errorf("create: %w", err)
	}

	return message, nil
}

func (r *messageRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *messageRepository) List(
	ctx context.Context,
) ([]*model.Message, error) {
	data := make([]*model.Message, listLimit)

	cursor, err := r.collection.Find(
		ctx,
		bson.D{},
		options.Find().SetLimit(listLimit),
	)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	if err := cursor.All(ctx, &data); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	return data, nil
}

func (r *messageRepository) NewMessage(
	ctx context.Context,
	delivery amqp091.Delivery,
) {
	msg, err := r.Create(ctx, string(delivery.Body))
	if err != nil {
		log.Printf("consume: failed to create message: %v\n", err)
		nack(delivery)

		return
	}

	if err := r.rabbitmqClient.Publisher.Publish(
		fmt.Sprintf(`
{
    "message": {
        "id": %q,
        "data": %q,
        "created": %q,
        "modified": %q
    },
    "provider": %q
}`,
			msg.ID.Hex(),
			msg.Data,
			msg.Created.Format(time.RFC3339),
			msg.Modified.Format(time.RFC3339),
			providerName,
		),
		amqp091.Table{
			"provider": strings.ToLower(providerName),
			"type":     rabbitmq.CreatedMessageKey,
		},
	); err != nil {
		log.Printf(
			"consume: failed to publish message creation message: %v\n",
			err,
		)
		nack(delivery)

		return
	}

	ack(delivery)
}

func ack(delivery amqp091.Delivery) {
	if err := delivery.Ack(false); err != nil {
		log.Printf("failed to ack delivery: %v\n", err)
	}
}

func nack(delivery amqp091.Delivery) {
	if err := delivery.Nack(false, true); err != nil {
		log.Printf("failed to nack delivery: %v\n", err)
	}
}
