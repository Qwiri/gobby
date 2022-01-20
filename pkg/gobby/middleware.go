package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (g *Gobby) Hook(app *fiber.App) error {
	app.Get(g.Prefix+"create", g.Router.routeLobbyCreate)
	app.Use(g.Prefix+"socket", g.Router.routeUpgradeWebsocket)
	app.Get(g.Prefix+"socket/:id", websocket.New(g.Router.routeGetSocket))
	return nil
}
