package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	api "github.com/penomatikus/golearning/go_websockets"
)

func main() {
	ws := newWsServer()
	ws.handle("/chat")
	ws.serve(":8080")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type (
	// A simple ID for connection identification
	connID int

	// Describes a unique connection to a wsSever
	wsConnection struct {
		id   connID
		conn *websocket.Conn
	}

	// Describes a websocket server with a pool of unique connections
	wsServer struct {
		connCount int
		pool      map[connID]*wsConnection
	}
)

func newWsServer() wsServer {
	log.Println("Hello Websockets!")
	return wsServer{
		pool: make(map[connID]*wsConnection),
	}
}

func (ws *wsServer) serve(port string) {
	log.Fatal(http.ListenAndServe(api.DefaultAddr, nil))
}

func (ws *wsServer) handle(route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("handleFunc: %s", err)
		}

		ws.connCount++
		wsConn := &wsConnection{
			id:   connID(ws.connCount),
			conn: conn,
		}

		log.Println("Client connected")
		ws.pool[wsConn.id] = wsConn
		wsConn.conn.WriteMessage(websocket.TextMessage, []byte("Schön, dass du da bist!"))

		err = ws.listen(wsConn)
		if err != nil {
			log.Printf("listen: %s", err)
			ws.updatePool(wsConn.id)
		}
	})
}

// listen on conn until an error occures
func (ws *wsServer) listen(wsConn *wsConnection) error {
	for {
		var m api.Message
		err := wsConn.conn.ReadJSON(&m)
		if err != nil {
			return err
		}

		switch m.Operation {
		case api.Broadcast:
			ws.broadcast(m.Message, wsConn.id)
		case api.Servertime:
			log.Printf("%s: Servertime is  %s", m.Username, time.Now().Format(time.RFC1123))
		default:
			log.Printf("%s: %s", m.Username, m.Message)
		}
	}
}

func (ws *wsServer) broadcast(broadcast string, broadcasterID connID) {
	for k := range ws.pool {
		if k != broadcasterID {
			err := ws.pool[k].conn.WriteMessage(websocket.TextMessage, []byte(broadcast))
			if err != nil {
				log.Fatalf("broadcast: %s", err)
			}
		}
	}
}

func (ws *wsServer) updatePool(id connID) {
	log.Printf("removed connection ID: %d\n", id)
	delete(ws.pool, id)
}
