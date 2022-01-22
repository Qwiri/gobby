package gobby

import (
	"errors"
	"github.com/gofiber/websocket/v2"
)

type BasicClientEvent func(*Client, *Lobby) error
type MessagedEvent func(*websocket.Conn, *Client, *Lobby, string) error
type ByteEvent func(*websocket.Conn, *Client, *Lobby, []byte) error

type EventType uint8

const (
	ReceiveEvent EventType = iota
	SendEvent
	JoinEvent
	LeaveEvent
)

func NewDispatcher(g *Gobby) *Dispatcher {
	return &Dispatcher{
		gobby:     g,
		listeners: make(map[EventType][]interface{}),
	}
}

var ErrInvalidListenerType = errors.New("invalid listener type")

func (d *Dispatcher) On(typ EventType, listener interface{}) error {
	// check if listener is valid
	if !IsListener(listener) {
		return ErrInvalidListenerType
	}
	d.listeners[typ] = append(d.listeners[typ], listener)
	return nil
}

func (d *Dispatcher) MustOn(typ EventType, listener interface{}) {
	if err := d.On(typ, listener); err != nil {
		panic(err)
	}
}

///

func (d *Dispatcher) Call(typ EventType, socket *websocket.Conn, client *Client, lobby *Lobby, message []byte) error {
	for _, ev := range d.listeners[typ] {
		switch f := ev.(type) {
		case BasicClientEvent:
			return f(client, lobby)
		case ByteEvent:
			return f(socket, client, lobby, message)
		case MessagedEvent:
			return f(socket, client, lobby, string(message))
		}
	}
	return nil
}

func (d *Dispatcher) handleMessage(socket *websocket.Conn, data []byte) {
	// TODO: implement this

	// TODO: ignore reply-messages (maybe add an extra event type for that)
}

func (d *Dispatcher) handleClose(socket *websocket.Conn) {
	// TODO: implement this
}
