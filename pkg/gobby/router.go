package gobby

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Router struct {
	g   *Gobby
	app *fiber.App
}

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

func (*Router) routeLobbyCreate(ctx *fiber.Ctx) error {
	return nil
}

func (*Router) routeGetSocket(socket *websocket.Conn) {
	gameID := socket.Params("id")
	log.Infof("[ws] got connection to id %s", gameID)

	// TODO: Send Lobby Version

	// TODO: make sure the lobby exists

	log.Infof("[ws] websocket connection with game %+v", gameID)
	for {
		if _, msg, err := socket.ReadMessage(); err != nil {
			_ = msg
			log.WithError(err).Warn("[ws] cannot read message from websocket")
			break
		}
		// TODO: handle message
	}

	// TODO: remove client from session
}

func (*Router) routeUpgradeWebsocket(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}
