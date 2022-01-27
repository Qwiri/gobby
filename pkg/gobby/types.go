package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type LobbyID string

type Client struct {
	Name   string
	Socket *websocket.Conn
	Role   Role
	Meta   interface{}
}

func NewClient(socket *websocket.Conn, name string) *Client {
	return &Client{
		Name:   name,
		Socket: socket,
		Role:   0,
	}
}

///

type Lobby struct {
	ID      LobbyID
	State   State
	Clients map[string]*Client
	Secret  string
	Meta    interface{}
}

func NewLobby(id LobbyID) *Lobby {
	return &Lobby{
		ID:      id,
		State:   0,
		Clients: make(map[string]*Client),
	}
}

type Dispatcher struct {
	gobby     *Gobby
	listeners map[EventType][]interface{}
}

type Router struct {
	g   *Gobby
	app *fiber.App
}
