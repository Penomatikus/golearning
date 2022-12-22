package main

import (
	"context"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
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

type chatServiceServerImpl struct {
	api.UnimplementedChatServiceServer
	registeredClients []string
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

func (impl *chatServiceServerImpl) RegisterClient(ctx context.Context, e *emptypb.Empty) (*emptypb.Empty, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return emptyMessage(), errors.New("no metadata found for register process")
	}

	token := uuid.New().String()
	trailer := metadata.New(map[string]string{"token": token})
	grpc.SetTrailer(ctx, trailer)

	impl.registeredClients = append(impl.registeredClients, token)
	return emptyMessage(), nil
}

func (impl *chatServiceServerImpl) WriteMessage(ctx context.Context, in *api.Message) (e *emptypb.Empty, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return emptyMessage(), errors.New("no metadata found for register process")
	}
	if md.Len() == 0 {
		return emptyMessage(), errors.New("empty metadata found for register process")
	}

	registered := false
	for _, s := range impl.registeredClients {
		if s == md.Get("token")[0] {
			registered = true
		}
	}

	if !registered {
		return emptyMessage(), errors.New("not registered")
	}

	return
}

func emptyMessage() *emptypb.Empty {
	empty := emptypb.Empty{}
	return &empty
}
