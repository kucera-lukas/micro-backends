package resolver

import (
	"github.com/99designs/gqlgen/graphql"

	"github.com/kucera-lukas/micro-backends/backend-service/gqlgen"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is a context struct.
type Resolver struct {
	config     *env.Config
	amqpClient *rabbitmq.Client
	controller controller.Controller
}

// NewSchema creates a new graphql.ExecutableSchema.
func NewSchema( //nolint:ireturn
	config *env.Config,
	amqpClient *rabbitmq.Client,
	controller controller.Controller,
) graphql.ExecutableSchema {
	return gqlgen.NewExecutableSchema(gqlgen.Config{
		Resolvers:  getResolver(config, amqpClient, controller),
		Directives: getDirective(),
		Complexity: getComplexity(),
	})
}

func getResolver(
	config *env.Config,
	amqpClient *rabbitmq.Client,
	controller controller.Controller,
) *Resolver {
	return &Resolver{
		config:     config,
		amqpClient: amqpClient,
		controller: controller,
	}
}

func getDirective() gqlgen.DirectiveRoot {
	return gqlgen.DirectiveRoot{}
}

func getComplexity() gqlgen.ComplexityRoot {
	return gqlgen.ComplexityRoot{} //nolint:exhaustivestruct
}
