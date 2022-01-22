package gobby

import "github.com/gofiber/websocket/v2"

func (g *Gobby) Broadcast(msg *Message) error {
	for _, l := range g.Lobbies {
		if err := l.Broadcast(msg); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gobby) BroadcastForce(msg *Message) {
	for _, l := range g.Lobbies {
		l.BroadcastForce(msg)
	}
}

func (l *Lobby) BroadcastForce(msg *Message) {
	d := msg.Marshal()
	for _, c := range l.Clients {
		_ = c.Socket.WriteMessage(websocket.TextMessage, d)
	}
}

func (l *Lobby) Broadcast(msg *Message) error {
	d := msg.Marshal()
	for _, c := range l.Clients {
		if err := c.Socket.WriteMessage(websocket.TextMessage, d); err != nil {
			return err
		}
	}
	return nil
}
