package gobby

import (
	"errors"
	"github.com/gofiber/websocket/v2"
)

var ErrInvalidListenerType = errors.New("invalid listener type")

type BasicEvent func(*Client, *Lobby) error
type MessageEvent func(*websocket.Conn, *Client, *Lobby, *Message) error

type EventType uint8

const (
	JoinEvent EventType = iota
	ReceiveEvent
	ReceiveReplyEvent
	LeaveEvent
)

func NewDispatcher(g *Gobby) *Dispatcher {
	return &Dispatcher{
		gobby:     g,
		listeners: make(map[EventType][]interface{}),
	}
}

// On registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, On returns an error
func (d *Dispatcher) On(typ EventType, listener ...interface{}) error {
	// check if listener is valid
	if !IsListener(listener) {
		return ErrInvalidListenerType
	}
	d.listeners[typ] = append(d.listeners[typ], listener...)
	return nil
}

// MustOn registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, the application panics
func (d *Dispatcher) MustOn(typ EventType, listener ...interface{}) {
	if err := d.On(typ, listener...); err != nil {
		panic(err)
	}
}

func (d *Dispatcher) Call(typ EventType, socket *websocket.Conn, client *Client, lobby *Lobby, message *Message) error {
	for _, ev := range d.listeners[typ] {
		switch f := ev.(type) {
		case BasicEvent:
			return f(client, lobby)
		case MessageEvent:
			return f(socket, client, lobby, message)
		}
	}
	return nil
}
