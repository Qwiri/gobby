package handlerutil

import "github.com/Qwiri/gobby/pkg/gobby"

func CreateListMessage(lobby *gobby.Lobby) *gobby.Message {
	names := make([]string, len(lobby.Clients))
	i := 0
	for _, c := range lobby.Clients {
		names[i] = c.Name
		i += 1
	}
	return gobby.NewBasicMessageWith("LIST", names...)
}
