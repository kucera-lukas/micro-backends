package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/adapter/repository"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/graphql"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/grpc"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/rabbitmq"
	"github.com/kucera-lukas/micro-backends/backend-service/pkg/infrastructure/router"
)

const (
	timeOutDeadline = time.Second * 30
	shutdownSignal  = 1
)

// Run runs the server with the given env.Config configuration.
func Run(config *env.Config) {
	srv := create(config)
	run(srv)
}

func create(config *env.Config) *http.Server {
	mongoClient := grpc.MustNewMongoClient(config)
	postgresClient := grpc.MustNewPostgresClient(config)
	rabbitmqClient := rabbitmq.MustNew(config.RabbitMQURI)

	ctrl := controller.Controller{
		Message: repository.NewMessageRepository(
			mongoClient,
			postgresClient,
			rabbitmqClient,
		),
	}

	gqlSrv := graphql.NewServer(config, rabbitmqClient, ctrl)
	ginRouter := router.New(config, gqlSrv)

	return &http.Server{ //nolint:exhaustivestruct
		Addr:         fmt.Sprintf(`0.0.0.0:%d`, config.Port),
		WriteTimeout: timeOutDeadline,
		ReadTimeout:  timeOutDeadline,
		IdleTimeout:  timeOutDeadline,
		Handler:      ginRouter,
	}
}

func run(srv *http.Server) {
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("listening on %s\n", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server terminated: %v\n", err)
		}
	}()

	channel := make(chan os.Signal, shutdownSignal)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(channel, os.Interrupt)

	// Block until we receive our signal.
	<-channel

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeOutDeadline)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("error shutting down the server: %v\n", err)
	}

	log.Println("server shutdown")
}
