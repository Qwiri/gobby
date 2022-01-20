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

type Gobby struct {
	Lobbies    map[LobbyID]*Lobby
	Prefix     string // default: /lobby/
	Dispatcher *Dispatcher
	Router     *Router
}

func New() (g *Gobby) {
	g = &Gobby{
		Lobbies: make(map[LobbyID]*Lobby),
		Prefix:  "/lobby/",
	}
	g.Dispatcher = &Dispatcher{g}
	g.Router = &Router{g}
	return
}

///

type Lobby struct {
	ID      LobbyID
	Clients map[string]*Client
}
