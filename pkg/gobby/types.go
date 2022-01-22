package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

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
	Prefix     string
	Dispatcher *Dispatcher
	Router     *Router
}

func New(app *fiber.App) (g *Gobby) {
	g = &Gobby{
		Lobbies: make(map[LobbyID]*Lobby),
		Prefix:  "/lobby/",
	}
	g.Dispatcher = NewDispatcher(g)
	// add router and register routes
	g.Router = NewRouter(g, app)
	g.Router.Hook()
	return
}

func (g *Gobby) OnReceive(events ...interface{}) {
	g.Dispatcher.onReceive = events
}

func (g *Gobby) OnSend(events ...interface{}) {
	g.Dispatcher.onSend = events
}

///

type Lobby struct {
	ID      LobbyID
	Clients map[string]*Client
}

///

type IncomingMessage struct {
	Client     *Client
	RequestUID string
}
