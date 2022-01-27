package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func NewRouter(g *Gobby, app *fiber.App) *Router {
	return &Router{
		g:   g,
		app: app,
	}
}

func (r *Router) Hook() {
	r.app.Get(r.g.Prefix+"create", r.routeLobbyCreate)
	r.app.Use(r.g.Prefix+"socket", r.routeUpgradeWebsocket)
	r.app.Get(r.g.Prefix+"socket/:id", websocket.New(r.routeGetSocket))
}
