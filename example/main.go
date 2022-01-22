package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()
	g := gobby.New(app)
	g.OnReceive(gobby.MessagedEvent(func(conn *websocket.Conn, lobby *gobby.Lobby, message string) error {
		log.Infof("Received message %s in lobby %s", message, lobby.ID)
		return nil
	}))

	g.BroadcastForce(websocket.TextMessage, []byte("HELLO"))
}
