# Go Websocket playground

An oversimplified chatroom implemented with websockets. Much space for improvement, however its intention is to get my hands on WebSockets with Go and get some learnings out of it.

## Key data: 
- 🏠 Home [ws://localhost:8080](ws://localhost:8080)
- 🔮 Uses [gorilla/websocket](github.com/gorilla/websocket)
- ✨ Go routines and channels 
- 🧙 Robust enough to not die on any client disconnect
- ✏️ Type `\t` client side to get the server time, which is useless since the sever tells the time for each message
- ✏️ Type `\b` client side to broadcast a message, which is useless since the sever displayes each message anyways
- ✏️ Type `\q` client side to logout
- 👋 Use any custom username 

## Start the server

_~/go_websockets/server:_ `go run .`

## Start a client

_~/go_websockets/client:_ `go run main.go -name Penomatikus`

## Example

### Broadcast a client:
```
Red: > Here is a nice video! https://www.youtube.com/watch?v=dQw4w9WgXcQ \b
```

View on other client: 
```
Server message: Here is a nice video! https://www.youtube.com/watch?v=dQw4w9WgXcQ
```
