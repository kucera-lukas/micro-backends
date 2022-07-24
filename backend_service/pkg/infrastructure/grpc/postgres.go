// nolint:ireturn
package grpc

import (
	"log"

	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/env"
	pbpostgres "github.com/kucera-lukas/micro-backends/backend-service/proto/postgres"
)

func MustNewPostgresClient(
	config *env.Config,
) pbpostgres.MessageServiceClient {
	client, err := NewPostgresClient(config)
	if err != nil {
		log.Panicf("grpc: failed to create postgres client: %v\n", err)
	}

	return client
}

func NewPostgresClient(
	config *env.Config,
) (pbpostgres.MessageServiceClient, error) {
	conn, err := dial(config.PostgresServiceAddress)
	if err != nil {
		return nil, err
	}

	return pbpostgres.NewMessageServiceClient(conn), nil
}
