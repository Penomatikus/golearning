package api

// defines a send message
type Message struct {
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Operation Operation `json:"operation"`
}

// defines certain server or client side behaviors e.g. broadcasting a message
type Operation string

const (
	// will broadcast the message
	Broadcast Operation = "\\b"
	// will print the current sever time
	Servertime Operation = "\\t"
	// will tell the sever about logging out
	Logout Operation = "\\q"
)

const DefaultHost string = "localhost:8080"
const DefaultAddr string = DefaultHost
const DefaultRoute string = "/chat"
