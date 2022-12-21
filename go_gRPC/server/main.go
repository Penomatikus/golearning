package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	api "github.com/penomatikus/golearning/go_gRPC/api"
)

func main() {
	gRPCchatService := newService()
	gRPCchatService.listenAndServe()
}

func newService() *chatServiceServerImpl {
	log.Println("Hello gRPC!")
	return &chatServiceServerImpl{
		registeredClients: make([]string, 0),
	}
}

func (impl *chatServiceServerImpl) listenAndServe() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterChatServiceServer(grpcServer, impl)
	log.Println("New Server started")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

type chatServiceServerImpl struct {
	api.UnimplementedChatServiceServer
	registeredClients []string
}

func (impl *chatServiceServerImpl) RegisterClient(ctx context.Context, e *emptypb.Empty) (*emptypb.Empty, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return impl.provideEmpty(), errors.New("no metadata found for register process")
	}

	token := time.Now().GoString()
	trailer := metadata.New(map[string]string{"token": token})
	grpc.SetTrailer(ctx, trailer)

	impl.registeredClients = append(impl.registeredClients, token)
	return impl.provideEmpty(), nil
}

func (impl *chatServiceServerImpl) provideEmpty() *emptypb.Empty {
	empty := emptypb.Empty{}
	return &empty
}
