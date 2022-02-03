package main

import (
	"fmt"
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/validate"
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

	g.Handle("CHAT", &gobby.Handler{
		Validation: validate.Schemes{
			validate.Number().As("time"),
			validate.String().Min(2).Max(16).As("username"),
			validate.String().Min(2).As("chat"),
		},
		Handler: func(event *gobby.Handle) error {
			username := event.Args.String("username")
			chat := event.Args.String("chat")

			fmt.Printf("%s wrote `%s` in lobby %s\n", username, chat, event.Lobby.ID)
			return nil
		},
	})

	g.MustOn(gobby.JoinEvent, gobby.AsyncBasicEvent(func(client *gobby.Client, lobby *gobby.Lobby) error {
		gobby.Infof(client, "joined lobby %s. Requesting client version ...", lobby.ID)

		// ask for client version and await response (blocks current goroutine)
		resp, err := gobby.NewBasicMessage("VERSION").SendAndAwaitReply(client.Socket, 10*time.Second)
		if err != nil {
			gobby.Warnf(client, "did not respond in time: %v", err)
			return err
		}

		gobby.Infof(client, "replied with version: %s", resp.Args[0].(string))
		return nil
	}))

	app.Listen(":8081")
}
