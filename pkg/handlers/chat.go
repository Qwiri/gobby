package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/validate"
	"github.com/apex/log"
)

//goland:noinspection GoUnusedGlobalVariable
var Chat = &gobby.Handler{
	Validation: validate.Schemes{
		validate.String().Min(1).Max(16).As("user"),
		validate.String().Min(1).Max(256).As("message"),
	},
	Handler: func(event *gobby.Handle) error {
		user := event.String("user")
		message := event.String("message")

		// build socket message and send to every client in lobby
		msg := gobby.NewBasicMessage("CHAT", user, message)

		var err error
		for _, c := range event.Lobby.Clients {
			if c == event.Client {
				err = event.Message.ReplyWith(c, *msg)
			} else {
				err = msg.SendTo(c)
			}
			if err != nil {
				log.WithError(err).Warnf("cannot send chat message to %s", c.Name)
			}
		}
		return nil
	},
}
