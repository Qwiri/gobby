package gobby

import (
	"github.com/gofiber/fiber/v2"
	"sync"
)

type Gobby struct {
	AppVersion string
	Lobbies    map[LobbyID]*Lobby
	lobbiesMu  sync.RWMutex
	Prefix     string
	// StrictLobby should be enabled if you want to have lobby created by calling the create endpoint
	// if StrictLobby is set to false, lobbies will be created by joining to a lobby
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
