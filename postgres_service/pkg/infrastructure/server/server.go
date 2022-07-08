package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/controller"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/adapter/repository"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/database"
	"github.com/kucera-lukas/micro-backends/postgres-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/postgres-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	timeOutDeadline = time.Second * 30
	shutdownSignal  = 1
)

type Server struct {
	proto.UnimplementedMessageServiceServer

	controller controller.Controller
}

func (s *Server) NewMessage(
	ctx context.Context,
	req *proto.NewMessageRequest,
) (*proto.NewMessageResponse, error) {
	id, err := s.controller.Message.Create(ctx, req.Data)
	if err != nil {
		return nil, err
	}

	return &proto.NewMessageResponse{Id: id}, nil
}

func (s *Server) MessageCount(
	ctx context.Context,
	req *proto.MessageCountRequest,
) (*proto.MessageCountResponse, error) {
	count, err := s.controller.Message.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.MessageCountResponse{Count: count}, nil
}

func (s *Server) GetMessages(
	ctx context.Context,
	req *proto.GetMessagesRequest,
) (*proto.GetMessagesResponse, error) {
	messageList, err := s.controller.Message.List(ctx)
	if err != nil {
		return nil, err
	}

	var messages []*proto.GetMessageResponse

	for _, msg := range messageList {
		messages = append(messages, &proto.GetMessageResponse{
			Id:       msg.Id,
			Data:     msg.Data,
			Created:  timestamppb.New(msg.Created),
			Modified: timestamppb.New(msg.Modified),
		})
	}

	return &proto.GetMessagesResponse{Messages: messages}, nil
}

// Run runs the server with the given env.Config configuration.
func Run(config *env.Config) {
	databaseClient := database.MustNew(config)
	defer databaseClient.Close()

	ctrl := controller.Controller{
		Message: repository.NewMessageRepository(databaseClient),
	}

	srv := grpc.NewServer()
	proto.RegisterMessageServiceServer(srv, &Server{controller: ctrl})
	reflection.Register(srv)

	address := fmt.Sprintf("0.0.0.0:%d", config.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("failed to listen to the address %s: %v", address, err)
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
