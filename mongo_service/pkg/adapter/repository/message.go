package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/model"
)

const (
	databaseName   = "micro_backends"
	collectionName = "messages"
)

// NewMessageRepository returns implementation of the controller.Message interface.
func NewMessageRepository(client *mongo.Client) controller.Message { //nolint:ireturn
	return &imageRepository{
		collection: client.Database(databaseName).Collection(collectionName),
	}
}

type imageRepository struct {
	collection *mongo.Collection
}

func (r *imageRepository) Create(
	ctx context.Context,
	data string,
) (primitive.ObjectID, error) {
	now := time.Now()

	message := &model.Message{
		ID:       primitive.NewObjectID(),
		Data:     data,
		Created:  now,
		Modified: now,
	}

	if _, err := r.collection.InsertOne(ctx, message); err != nil {
		return message.ID, err
	}

	return message.ID, nil
}

func (r *imageRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *imageRepository) List(ctx context.Context) ([]model.Message, error) {
	var data []model.Message

	cursor, err := r.collection.Find(
		ctx,
		bson.D{},
		options.Find().SetLimit(100),
	)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, nil
}
