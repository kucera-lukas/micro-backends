package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/model"
)

// Message Controller interface.
type Message interface {
	Create(ctx context.Context, data string) (primitive.ObjectID, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]model.Message, error)
}
