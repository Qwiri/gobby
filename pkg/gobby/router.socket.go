package gobby

import (
	"errors"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var ErrLobbyNotFound = errors.New("lobby not found")

func (*Router) routeUpgradeWebsocket(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (r *Router) routeGetSocket(socket *websocket.Conn) {
	lobbyID := LobbyID(socket.Params("id"))
	log.Infof("[ws] got new connection requesting lobby %s", lobbyID)

	// send gobby and app version to connected client
	if err := NewBasicMessage("VERSION", Version, r.g.AppVersion).Send(socket); err != nil {
		log.Warnf("failed to send version to %v: %v", socket.RemoteAddr(), err)
		return
	}

	// check if lobby exists
	r.g.lobbiesMu.RLock()
	defer r.g.lobbiesMu.RUnlock()

	lobby, ok := r.g.Lobbies[lobbyID]
	if !ok {
		_ = NewErrorMessage(ErrLobbyNotFound).Send(socket)
		return
	}
	log.Infof("[ws] requested lobby %s exists. Entering message loop.", lobbyID)

	for {
		if mt, msg, err := socket.ReadMessage(); err != nil {
			if err == websocket.ErrCloseSent {
				Infof(socket, "closed the connection")
			} else {
				Warnf(socket, "cannot read message: %s", err.Error())
			}
			break
		} else {
			if mt != websocket.TextMessage {
				Warnf(socket, "sent non-text data")
				continue
			}
			if err = r.g.Dispatcher.handleMessage(socket, lobby, msg); err != nil {
				if err = NewErrorMessage(err).Send(socket); err != nil {
					Warnf(socket, "cannot send panic message: %v", err)
				}
			}
		}
	}

	r.g.Dispatcher.handleClose(socket)
}
