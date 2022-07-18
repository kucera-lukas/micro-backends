package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/resolver"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
)

const (
	lruQueryCacheSize = 1000
	complexityLimit   = 1000
)

// NewServer generates a new handler.Server.
func NewServer(
	config *env.Config,
	rabbitmqClient *rabbitmq.Client,
	controller controller.Controller,
) *handler.Server {
	srv := handler.New(resolver.NewSchema(config, rabbitmqClient, controller))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{})
	srv.SetQueryCache(lru.New(lruQueryCacheSize))
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(complexityLimit))
	//srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
	//	return util.NewInternalServerError(ctx, fmt.Sprintf(`%v`, err))
	//})

	return srv
}
