package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"time"
)

func main() {
	app := fiber.New()
	g := gobby.New(app)
	d := g.Dispatcher

	d.MustOn(gobby.JoinEvent, gobby.BasicClientEvent(func(client *gobby.Client, lobby *gobby.Lobby) error {
		gobby.Infof(client, "joined lobby %s. Requesting client version ...", lobby.ID)

		// ask for client version
		resp, err := gobby.NewBasicMessage("VERSION").SendAndAwaitReply(client.Socket, 2*time.Second)
		if err != nil {
			log.Warnf("[%s] did not respond in time: %v", client.Name, err)
			return err
		}

		gobby.Infof(client, "got version: %s", resp.Args[0].(string))
		return nil
	}))
}
