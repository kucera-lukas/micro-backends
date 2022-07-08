package controller

import (
	"context"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/model"
)

// Message Controller interface.
type Message interface {
	Create(ctx context.Context, data string) (uint32, error)
	Count(ctx context.Context) (uint32, error)
	List(ctx context.Context) ([]model.Message, error)
}
