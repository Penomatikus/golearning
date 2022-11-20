package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	api "github.com/penomatikus/golearning/go_websockets"
)

func main() {
	name := flag.String("name", "user", "A username to display")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: api.DefaultHost, Path: api.DefaultRoute}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dail:", err)
	}
	defer conn.Close()
	go updates(conn, *name)

	messagechan := make(chan api.Message)
	closechan := make(chan struct{})
	go chat(*name, messagechan, closechan)

	for {
		select {
		case m := <-messagechan:
			err := conn.WriteJSON(m)
			if err != nil {
				log.Println("write: ", err)
				return
			}
		case <-closechan:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye!"))
			if err != nil {
				log.Println("on close: ", err)
			}
			return
		}

	}
}

func chat(username string, m chan api.Message, c chan<- struct{}) {
	for {
		fmt.Printf("%s: > ", username)

		r := bufio.NewReader(os.Stdin)
		in, _ := r.ReadString('\n')
		chatData := api.Message{
			Username: username,
			Message:  strings.TrimSuffix(in, "\n"),
		}

		if strings.Contains(in, string(api.Servertime)) {
			chatData.Operation = api.Servertime
		}

		if strings.Contains(in, string(api.Broadcast)) {
			chatData.Operation = api.Broadcast
		}

		if strings.Contains(in, string(api.Logout)) {
			c <- struct{}{}
			return
		}
		m <- chatData
	}
}

func updates(conn *websocket.Conn, username string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		fmt.Printf("\nServer message: %s\n%s: >", message, username)
	}
}
