package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/repository"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/database"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/rabbitmq"
	"github.com/kucera-lukas/micro-backends/postgres-service/proto"
)

const (
	timeOutDeadline = time.Second * 30
	shutdownSignal  = 1
)

type Server struct {
	pbpostgres.UnimplementedMessageServiceServer

	controller controller.Controller
}

func (s *Server) NewMessage(
	ctx context.Context,
	req *pbpostgres.NewMessageRequest,
) (*pbpostgres.NewMessageResponse, error) {
	message, err := s.controller.Message.Create(ctx, req.Data)
	if err != nil {
		return nil, err
	}

	return &pbpostgres.NewMessageResponse{
		Id:       strconv.Itoa(message.Id),
		Data:     message.Data,
		Created:  timestamppb.New(message.Created),
		Modified: timestamppb.New(message.Modified),
	}, nil
}

func (s *Server) MessageCount(
	ctx context.Context,
	req *pbpostgres.MessageCountRequest,
) (*pbpostgres.MessageCountResponse, error) {
	count, err := s.controller.Message.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &pbpostgres.MessageCountResponse{Count: count}, nil
}

func (s *Server) GetMessage(
	ctx context.Context,
	req *pbpostgres.GetMessageRequest,
) (*pbpostgres.GetMessageResponse, error) {
	message, err := s.controller.Message.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pbpostgres.GetMessageResponse{
		Id:       strconv.Itoa(message.Id),
		Data:     message.Data,
		Created:  timestamppb.New(message.Created),
		Modified: timestamppb.New(message.Modified),
	}, nil
}

func (s *Server) GetMessages(
	ctx context.Context,
	req *pbpostgres.GetMessagesRequest,
) (*pbpostgres.GetMessagesResponse, error) {
	messageList, err := s.controller.Message.List(ctx)
	if err != nil {
		return nil, err
	}

	var messages []*pbpostgres.GetMessageResponse

	for _, msg := range messageList {
		messages = append(messages, &pbpostgres.GetMessageResponse{
			Id:       strconv.Itoa(msg.Id),
			Data:     msg.Data,
			Created:  timestamppb.New(msg.Created),
			Modified: timestamppb.New(msg.Modified),
		})
	}

	return &pbpostgres.GetMessagesResponse{Messages: messages}, nil
}

// Run runs the server with the given env.Config configuration.
func Run(config *env.Config) {
	pgxPool := database.MustNew(config)
	defer pgxPool.Close()
	rabbitmqClient := rabbitmq.MustNew(config.RabbitMQURI)
	defer rabbitmqClient.Close()

	ctrl := controller.Controller{
		Message: repository.NewMessageRepository(pgxPool, rabbitmqClient),
	}
	if err := ctrl.Setup(rabbitmqClient); err != nil {
		log.Panicf("failed to setup controller: %v\n", err)
	}

	srv := grpc.NewServer()
	pbpostgres.RegisterMessageServiceServer(srv, &Server{controller: ctrl})
	reflection.Register(srv)

	address := fmt.Sprintf("0.0.0.0:%d", config.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf(
			"failed to listen to the address %s: %v\n",
			address,
			err,
		)
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("listening on %s", lis.Addr().String())

		if err := srv.Serve(lis); err != nil {
			log.Printf("server terminated: %v", err)
		}
	}()

	channel := make(chan os.Signal, shutdownSignal)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(channel, os.Interrupt)

	// Block until we receive our signal.
	<-channel

	// Create a deadline to wait for.
	_, cancel := context.WithTimeout(context.Background(), timeOutDeadline)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.GracefulStop()

	log.Println("server gracefully stopped")
}
