package handlers

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/validate"
)

//goland:noinspection GoUnusedGlobalVariable
var Chat = &gobby.Handler{
	Validation: validate.Schemes{
		validate.String().Min(1).Max(512).As("message"),
	},
	Handler: func(event *gobby.Handle) error {
		message := event.String("message")

		// build socket message and send to every client in lobby
		msg := gobby.NewBasicMessage("CHAT", event.Client.Name, message)
		event.Lobby.BroadcastForce(msg)

		return nil
	},
}
