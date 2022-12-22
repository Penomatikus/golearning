package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	api "github.com/penomatikus/golearning/go_gRPC/api"
)

func main() {
	name := flag.String("name", "user", "A username to display")
	flag.Parse()

	ctx := context.Background()

	chatClient := newChatClient(*name)
	conn := chatClient.dailAndConnect(ctx)
	defer conn.Close()

	chatClient.registerWithContext(ctx)
	print(chatClient.serverAuthToken)

}

type chatClient struct {
	api.ChatServiceClient
	name            string
	serverAuthToken serverAuthToken
}

func newChatClient(name string) *chatClient {
	return &chatClient{
		name: name,
	}
}

func (cc *chatClient) dailAndConnect(ctx context.Context) *grpc.ClientConn {
	defaultHost := ":9000"
	defaultDailOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, defaultHost, defaultDailOpt)
	if err != nil {
		log.Fatalf("Could not connect to server: %s", err)
	}
	cc.ChatServiceClient = api.NewChatServiceClient(conn)
	return conn
}

type serverAuthToken string

// registerWithContext "fakes" a register process to the server for later a later "authtoken"
// The server receives metadata which content is irrelevant but must exists and sends a UUID in
// the trailer back. This is just for playing around with metadata and trailers
func (cc *chatClient) registerWithContext(ctx context.Context) {
	md := metadata.New(map[string]string{"name": cc.name})
	ctx = metadata.NewOutgoingContext(ctx, md)

	var trailer metadata.MD
	_, err := cc.RegisterClient(
		ctx,
		emptyMessage(),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		log.Fatalf("Error in RegisterClient(): %s", err)
	}

	if len(trailer.Get("token")) == 0 {
		log.Fatal("Error in RegisterClienet(): The sever didn't send a token.")
	}

	cc.serverAuthToken = serverAuthToken(trailer.Get("token")[0])
}

func emptyMessage() *emptypb.Empty {
	empty := emptypb.Empty{}
	return &empty
}

// func chat(username string, m chan api.Message, c chan<- struct{}) {
// 	for {
// 		fmt.Printf("%s: > ", username)

// 		r := bufio.NewReader(os.Stdin)
// 		in, _ := r.ReadString('\n')
// 		chatData := api.Message{
// 			Username: username,
// 			Message:  strings.TrimSuffix(in, "\n"),
// 		}

// 		if strings.Contains(in, string(api.Servertime)) {
// 			chatData.Operation = api.Servertime
// 		}

// 		if strings.Contains(in, string(api.Broadcast)) {
// 			chatData.Operation = api.Broadcast
// 		}

// 		if strings.Contains(in, string(api.Logout)) {
// 			c <- struct{}{}
// 			return
// 		}
// 		m <- chatData
// 	}
// }
