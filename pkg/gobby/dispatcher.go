package gobby

import "github.com/gofiber/websocket/v2"

type BasicEvent func(*websocket.Conn, *Lobby) error
type MessagedEvent func(*websocket.Conn, *Lobby, string) error
type ByteEvent func(*websocket.Conn, *Lobby, []byte) error

type Dispatcher struct {
	gobby     *Gobby
	onReceive []interface{}
	onSend    []interface{}
}

///

func NewDispatcher(g *Gobby) *Dispatcher {
	return &Dispatcher{
		gobby:     g,
		onReceive: make([]interface{}, 0),
		onSend:    make([]interface{}, 0),
	}
}

func (d *Dispatcher) AddReceive(events ...interface{}) {
	d.onReceive = append(d.onReceive, events...)
}

func (d *Dispatcher) AddSend(events ...interface{}) {
	d.onSend = append(d.onSend, events...)
}

///

func (d *Dispatcher) call(events []interface{}, socket *websocket.Conn, lobby *Lobby, message []byte) error {
	for _, ev := range events {
		switch f := ev.(type) {
		case BasicEvent:
			return f(socket, lobby)
		case ByteEvent:
			return f(socket, lobby, message)
		case MessagedEvent:
			return f(socket, lobby, string(message))
		}
	}
	return nil
}

func (d *Dispatcher) CallMessageReceive(socket *websocket.Conn, lobby *Lobby, message []byte) error {
	return d.call(d.onReceive, socket, lobby, message)
}

func (d *Dispatcher) CallMessageSend(socket *websocket.Conn, lobby *Lobby, message []byte) error {
	return d.call(d.onSend, socket, lobby, message)
}
