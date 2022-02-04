package gobby

import (
	"github.com/gofiber/fiber/v2"
)

func (g *Gobby) LobbyExists(id LobbyID) (ok bool) {
	g.lobbiesMu.RLock()
	_, ok = g.Lobbies[id]
	g.lobbiesMu.RUnlock()
	return
}

func (r *Router) routeLobbyCreate(ctx *fiber.Ctx) error {
	// find a free lobby ID
	var id LobbyID
	for id == "" || r.g.LobbyExists(id) {
		id = LobbyID(generateRandomString(8))
	}

	// create lobby
	lobby := NewLobby(id)

	// dispatch lobby create event
	r.g.Dispatcher.call(lobbyCreateType, &LobbyCreate{
		Lobby: lobby,
		Addr:  ctx.IP(),
	})

	r.g.lobbiesMu.Lock()
	r.g.Lobbies[id] = lobby
	r.g.lobbiesMu.Unlock()

	return ctx.JSON(lobby)
}
