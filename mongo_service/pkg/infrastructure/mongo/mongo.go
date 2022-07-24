package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/env"
)

// MustNew ensure that a new mongo.Client is created and panics if not.
func MustNew(config *env.Config) *mongo.Client {
	client, err := New(context.Background(), config)
	if err != nil {
		log.Panic(err)
	}

	return client
}

// New tries to create a new mongo.Client, returning error if unsuccessful.
func New(ctx context.Context, config *env.Config) (*mongo.Client, error) {
	credential := options.Credential{ // nolint:exhaustivestruct
		Username: config.MongoDBUsername,
		Password: config.MongoDBPassword,
	}

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(config.MongoDBURI).SetAuth(credential),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client, nil
}

func Disconnect(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Panicf("failed to disconnect from mongodb: %v", err)
	}
}
