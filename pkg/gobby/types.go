package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type LobbyID string

type Client struct {
	Name   string
	Socket *websocket.Conn
	Meta   interface{}
}

///

type Lobby struct {
	ID      LobbyID
	Clients map[string]*Client
	Meta    interface{}
}

func NewLobby(id LobbyID) *Lobby {
	return &Lobby{
		ID:      id,
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
