package gobby

import "github.com/gofiber/websocket/v2"

///

type (
	LobbyID    string
	ClientName string
)

///

type Client struct {
	Name   ClientName
	Socket *websocket.Conn
	Meta   interface{}
}

///

type Lobby struct {
	ID      LobbyID
	Prefix  string // default: /lobby/
	Version string
}
