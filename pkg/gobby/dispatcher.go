package gobby

import (
	"errors"
)

var ErrInvalidListenerType = errors.New("invalid listener type")

type eventType uint8

type Dispatcher struct {
	gobby    *Gobby
	handlers map[eventType][]*eventHandlerWrapper
}

func NewDispatcher(g *Gobby) *Dispatcher {
	return &Dispatcher{
		gobby:    g,
		handlers: make(map[eventType][]*eventHandlerWrapper),
	}
}

type eventHandlerWrapper struct {
	handler eventHandler
}

type eventHandler interface {
	Type() eventType
	Handle(i interface{})
}

// On registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, On returns an error
func (d *Dispatcher) On(listener ...interface{}) error {
	// check if listener is valid
	for _, l := range listener {
		// find event handler for listener
		eh := d.getEventHandler(l)
		if eh == nil {
			return ErrInvalidListenerType
		}
		d.addEventHandler(eh)
	}
	return nil
}

// MustOn registers a BasicClientEvent, ByteEvent or MessagedEvent
// if the listener argument is no valid listener, the application panics
func (d *Dispatcher) MustOn(listener ...interface{}) {
	if err := d.On(listener...); err != nil {
		panic(err)
	}
}

func (d *Dispatcher) addEventHandler(eventHandler eventHandler) {
	wrap := &eventHandlerWrapper{eventHandler}
	d.handlers[eventHandler.Type()] = append(d.handlers[eventHandler.Type()], wrap)
}

func (d *Dispatcher) call(typ eventType, i interface{}) {
	for _, h := range d.handlers[typ] {
		h.handler.Handle(i)
	}
	return
}
