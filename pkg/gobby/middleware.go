package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (l *Lobby) Hook(app *fiber.App) error {
	r := &Router{l}
	app.Get(l.Prefix+"create", r.routeLobbyCreate)
	app.Use(l.Prefix+"socket", r.routeUpgradeWebsocket)
	app.Get(l.Prefix+"socket/:id", websocket.New(r.routeGetSocket))
	return nil
}
