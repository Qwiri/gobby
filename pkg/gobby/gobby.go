package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"strings"
	"sync"
)

type Gobby struct {
	AppVersion string
	Lobbies    map[LobbyID]*Lobby
	lobbiesMu  sync.RWMutex
	Prefix     string
	Dispatcher *Dispatcher
	Router     *Router
	Handlers   map[string]*Handler
}

func New(app *fiber.App) (g *Gobby) {
	g = &Gobby{
		Lobbies:  make(map[LobbyID]*Lobby),
		Handlers: make(map[string]*Handler),
		Prefix:   "/lobby/",
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
	delete(lobby.Clients, strings.ToLower(client.Name))
	// call leave event
	g.Dispatcher.call(leaveType, &Leave{
		Client: client,
		Lobby:  lobby,
	})
	// send PLAYER_LEAVE to all players
	lobby.BroadcastForce(NewBasicMessageWith[string]("PLAYER_LEAVE", client.Name))
	// send LIST message to all clients
	lobby.BroadcastForce(CreateListMessage(lobby))
}

// Aliases

// On registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, On returns an error
func (g *Gobby) On(listener ...interface{}) error {
	return g.Dispatcher.On(listener...)
}

// MustOn registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, the application panics
func (g *Gobby) MustOn(listener ...interface{}) {
	g.Dispatcher.MustOn(listener...)
}

type Handlers map[string]*Handler

func (g *Gobby) Handle(name string, handler *Handler) {
	g.Handlers[strings.ToLower(name)] = handler
}

func (g *Gobby) HandleAll(h Handlers) {
	for n, v := range h {
		g.Handle(n, v)
	}
}
