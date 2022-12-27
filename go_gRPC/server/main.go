package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	api "github.com/penomatikus/golearning/go_gRPC/shared/api"
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

type broker struct {
}

func (b *broker) broadcast() {

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

func (impl *chatServiceServerImpl) PublishMessage(stream api.ChatService_PublishMessageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Print(in.GetMessage())
	}
}

func emptyMessage() *emptypb.Empty {
	empty := emptypb.Empty{}
	return &empty
}
