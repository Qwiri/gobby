package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/validate"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
		Handler: gobby.MessagedHandler(func(conn *websocket.Conn, lobby *gobby.Lobby, client *gobby.Client, message *gobby.Message) error {
			return nil
		}),
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
