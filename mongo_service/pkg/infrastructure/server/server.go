package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/adapter/repository"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/mongo"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/rabbitmq"
	pbmongo "github.com/kucera-lukas/micro-backends/mongo-service/proto"
)

const (
	timeOutDeadline = time.Second * 30
	shutdownSignal  = 1
)

type Server struct {
	pbmongo.UnimplementedMessageServiceServer

	controller controller.Controller
}

func (s *Server) NewMessage(
	ctx context.Context,
	req *pbmongo.NewMessageRequest,
) (*pbmongo.NewMessageResponse, error) {
	message, err := s.controller.Message.Create(ctx, req.Data)
	if err != nil {
		return nil, fmt.Errorf("new_message: %w", err)
	}

	return &pbmongo.NewMessageResponse{
		Id:       message.ID.Hex(),
		Data:     message.Data,
		Created:  timestamppb.New(message.Created),
		Modified: timestamppb.New(message.Modified),
	}, nil
}

func (s *Server) MessageCount(
	ctx context.Context,
	req *pbmongo.MessageCountRequest,
) (*pbmongo.MessageCountResponse, error) {
	count, err := s.controller.Message.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("message_count: %w", err)
	}

	return &pbmongo.MessageCountResponse{Count: count}, nil
}

func (s *Server) GetMessage(
	ctx context.Context,
	req *pbmongo.GetMessageRequest,
) (*pbmongo.GetMessageResponse, error) {
	message, err := s.controller.Message.Get(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("get_message: %w", err)
	}

	return &pbmongo.GetMessageResponse{
		Id:       message.ID.Hex(),
		Data:     message.Data,
		Created:  timestamppb.New(message.Created),
		Modified: timestamppb.New(message.Modified),
	}, nil
}

func (s *Server) GetMessages(
	ctx context.Context,
	req *pbmongo.GetMessagesRequest,
) (*pbmongo.GetMessagesResponse, error) {
	messageList, err := s.controller.Message.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get_messages: %w", err)
	}

	messages := make([]*pbmongo.GetMessageResponse, len(messageList))

	for i, msg := range messageList {
		messages[i] = &pbmongo.GetMessageResponse{
			Id:       msg.ID.Hex(),
			Data:     msg.Data,
			Created:  timestamppb.New(msg.Created),
			Modified: timestamppb.New(msg.Modified),
		}
	}

	return &pbmongo.GetMessagesResponse{Messages: messages}, nil
}

// Run runs the server with the given env.Config configuration.
func Run(config *env.Config) {
	rabbitmqClient := rabbitmq.MustNew(config.RabbitMQURI)
	defer rabbitmqClient.Close()

	mongoClient := mongo.MustNew(config)
	defer mongo.Disconnect(mongoClient)

	ctrl := controller.Controller{
		Message: repository.NewMessageRepository(mongoClient, rabbitmqClient),
	}
	if err := ctrl.Setup(rabbitmqClient); err != nil {
		log.Panicf("failed to setup controller: %v\n", err)
	}

	srv := grpc.NewServer()
	pbmongo.RegisterMessageServiceServer(
		srv,
		&Server{controller: ctrl}, // nolint:exhaustivestruct
	)
	reflection.Register(srv)

	address := fmt.Sprintf("0.0.0.0:%d", config.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf(
			"failed to listen to the address %s: %v",
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
