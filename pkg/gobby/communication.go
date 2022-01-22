package gobby

import "github.com/gofiber/websocket/v2"

/// Gobby

func (g *Gobby) Broadcast(messageType int, message []byte) error {
	for _, l := range g.Lobbies {
		if err := l.Broadcast(messageType, message); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gobby) BroadcastForce(messageType int, message []byte) {
	for _, l := range g.Lobbies {
		l.BroadcastForce(messageType, message)
	}
}

//// Lobby

// Any

func (l *Lobby) BroadcastForce(messageType int, message []byte) {
	for _, c := range l.Clients {
		_ = c.Socket.WriteMessage(messageType, message)
	}
}

func (l *Lobby) Broadcast(messageType int, message []byte) error {
	for _, c := range l.Clients {
		if err := c.Socket.WriteMessage(messageType, message); err != nil {
			return err
		}
	}
	return nil
}

// Bytes

func (l *Lobby) BroadcastBytesForce(message []byte) {
	l.BroadcastForce(websocket.BinaryMessage, message)
}

func (l *Lobby) BroadcastBytes(message []byte) error {
	return l.Broadcast(websocket.BinaryMessage, message)
}

// String

func (l *Lobby) BroadcastStringForce(message string) {
	l.BroadcastForce(websocket.TextMessage, []byte(message))
}

func (l *Lobby) BroadcastString(message string) error {
	return l.Broadcast(websocket.TextMessage, []byte(message))
}
