package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	api "github.com/penomatikus/golearning/go_gRPC/api"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterChatServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

type server struct {
	api.UnimplementedChatServiceServer
}

func (s *server) SayHello(ctx context.Context, in *api.Message) (*api.Message, error) {
	log.Printf("Received: %v", in.GetBody())
	return &api.Message{Body: "Hallo" + in.GetBody()}, nil

}
