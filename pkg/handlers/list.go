package handlers

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/validate"
)

//goland:noinspection GoUnusedGlobalVariable
var List = &gobby.Handler{
	Validation: validate.Schemes{},
	Handler: func(event *gobby.Handle) error {
		return event.Message.ReplyWith(event.Client, *gobby.CreateListMessage(event.Lobby))
	},
}
