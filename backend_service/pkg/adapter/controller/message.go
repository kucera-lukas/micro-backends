package controller

import (
	"context"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
)

// Message Controller interface.
type Message interface {
	Get(
		ctx context.Context,
		id string,
		provider gqlgen.MessageProvider,
	) (*gqlgen.Message, error)
	List(
		ctx context.Context,
		provider gqlgen.MessageProvider,
	) ([]*gqlgen.Message, error)
	ListAll(
		ctx context.Context,
	) ([]*gqlgen.Message, error)
	Count(
		ctx context.Context,
		provider gqlgen.MessageProvider,
	) (int64, error)
	CountAll(
		ctx context.Context,
	) (int64, error)
	Create(
		ctx context.Context,
		data string,
		provider gqlgen.MessageProvider,
	) (*gqlgen.Message, error)
	CreateAll(
		ctx context.Context,
		data string,
	) (string, error)
}
