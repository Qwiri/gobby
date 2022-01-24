package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type Gobby struct {
	AppVersion string
	Lobbies    map[LobbyID]*Lobby
	lobbiesMu  sync.RWMutex
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

func (g *Gobby) BySocket(socket *websocket.Conn) (*Lobby, *Client, bool) {
	for _, l := range g.Lobbies {
		for _, c := range l.Clients {
			if c.Socket == socket {
				return l, c, true
			}
		}
	}
	return nil, nil, false
}

func (g *Gobby) RemoveClient(lobby *Lobby, client *Client) {
	delete(lobby.Clients, client.Name)
	if err := g.Dispatcher.Call(LeaveEvent, client.Socket, client, lobby, nil); err != nil {
		Warnf(client, "cannot call leave event: %v", err)
	}
}

// Aliases

// On registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, On returns an error
func (g *Gobby) On(typ EventType, listener ...interface{}) error {
	return g.Dispatcher.On(typ, listener...)
}

// MustOn registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, the application panics
func (g *Gobby) MustOn(typ EventType, listener ...interface{}) {
	g.Dispatcher.MustOn(typ, listener...)
}
