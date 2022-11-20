module github.com/penomatikus/golearning/go_websockets/server

go 1.19

require (
	github.com/gorilla/websocket v1.5.0
	github.com/penomatikus/golearning/go_websockets v0.0.0-00010101000000-000000000000
)

replace github.com/penomatikus/golearning/go_websockets => ./../
