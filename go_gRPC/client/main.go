package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	api "github.com/penomatikus/golearning/go_gRPC/api"
)

const defaultHost = ":9000"

var /* const */ defaultDailOpt = grpc.WithTransportCredentials(insecure.NewCredentials())

func main() {
	name := flag.String("name", "user", "A username to display")
	flag.Parse()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, defaultHost, defaultDailOpt)
	if err != nil {
		panic("Could not connect to server")
	}
	defer conn.Close()

	// messagechan := make(chan api.Message)
	// closechan := make(chan struct{})

	client := api.NewChatServiceClient(conn)
	// body := fmt.Sprintf("Hello From Client \"%s\"!", *name)

	// response, err := client.SayHello(context.Background(), &api.Message{Body: body})
	md := metadata.New(map[string]string{"name": *name})
	ctx = metadata.NewOutgoingContext(ctx, md)

	empty := emptypb.Empty{}
	var header, trailer metadata.MD
	_, err = client.RegisterClient(
		ctx,
		&empty,
		grpc.Header(&header),   // will retrieve header
		grpc.Trailer(&trailer), // will retrieve trailer)
	)
	if err != nil {
		log.Fatalf("Error in RegisterClient(): %s", err)
	}

	fmt.Printf("header: %v\n", header)
	fmt.Printf("trailer: %v", trailer)

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
