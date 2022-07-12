package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/handlers"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix() << 1)
	log.SetLevel(log.DebugLevel)
}

func main() {

	app := fiber.New()
	g := gobby.New(app)

	g.HandleAll(gobby.Handlers{
		"CHAT": handlers.Chat,
		"LIST": handlers.List,
	})

	g.MustOn(func(event *gobby.PreJoinEvent) {
		client, lobby := event.Client, event.Lobby
		gobby.Infof(client, "joined lobby %s. Requesting client version.", lobby.ID)

		go func() {
			// ask for client version and await response (blocks current goroutine)
			resp, err := gobby.NewBasicMessage("VERSION").SendAndAwaitReply(client.Socket, 10*time.Second)
			if err != nil {
				gobby.Warnf(client, "did not respond in time: %v", err)
			} else {
				gobby.Infof(client, "replied with version: %s", resp.Args[0].(string))
			}
		}()
	})

	if err := app.Listen(":8081"); err != nil {
		log.WithError(err).Warn("cannot serve")
	}
}
