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
)

// <--
// oversimplified contract
type (
	operation string
	message   struct {
		Username  string
		Message   string
		Operation operation
	}
)

const (
	// will broadcast the message
	Broadcast operation = "\\b"
	// will print the current sever time
	Servertime operation = "\\t"
	// will tell the sever about logging out
	Logout operation = "\\q"
)

// -->

func main() {
	name := flag.String("name", "user", "A username to display")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/chat"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dail:", err)
	}
	defer conn.Close()
	go updates(conn, *name)

	messagechan := make(chan message)
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

func chat(username string, m chan message, c chan<- struct{}) {
	for {
		fmt.Printf("%s: > ", username)

		r := bufio.NewReader(os.Stdin)
		in, _ := r.ReadString('\n')
		chatData := message{
			Username: username,
			Message:  strings.TrimSuffix(in, "\n"),
		}

		if strings.Contains(in, string(Servertime)) {
			chatData.Operation = Servertime
		}

		if strings.Contains(in, string(Broadcast)) {
			chatData.Operation = Broadcast
		}

		if strings.Contains(in, string(Logout)) {
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
