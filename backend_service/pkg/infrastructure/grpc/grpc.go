package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func dial(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(target, getOpts()...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}
