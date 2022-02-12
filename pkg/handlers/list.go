package handlers

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/util"
	"github.com/Qwiri/gobby/pkg/validate"
)

//goland:noinspection GoUnusedGlobalVariable
var List = &gobby.Handler{
	Validation: validate.Schemes{},
	Handler: func(event *gobby.Handle) error {
		names := make([]string, len(event.Lobby.Clients))
		i := 0
		for _, c := range event.Lobby.Clients {
			names[i] = c.Name
			i += 1
		}
		return event.Message.ReplyWith(event.Client,
			*gobby.NewBasicMessage("LIST", util.StringToAnyArray(names)...))
	},
}
