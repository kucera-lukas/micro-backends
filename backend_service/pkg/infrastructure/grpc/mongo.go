// nolint:ireturn
package grpc

import (
	"log"

	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/env"
	pbmongo "github.com/kucera-lukas/micro-backends/backend-service/proto/mongo"
)

func MustNewMongoClient(
	config *env.Config,
) pbmongo.MessageServiceClient {
	client, err := NewMongoClient(config)
	if err != nil {
		log.Panicf("grpc: failed to create mongo client: %v\n", err)
	}

	return client
}

func NewMongoClient(
	config *env.Config,
) (pbmongo.MessageServiceClient, error) {
	conn, err := dial(config.MongoServiceAddress)
	if err != nil {
		return nil, err
	}

	return pbmongo.NewMessageServiceClient(conn), nil
}
