package gobby

func CreateListMessage(lobby *Lobby) *Message {
	names := make([]string, len(lobby.Clients))
	i := 0
	for _, c := range lobby.Clients {
		names[i] = c.Name
		i += 1
	}
	return NewBasicMessageWith("LIST", names...)
}
