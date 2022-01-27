package gobby

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/websocket/v2"
	"strings"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
	ErrPayloadTooLarge = errors.New("payload too large")

	ErrArgs              = errors.New("args mismatch")
	ErrSecretMismatch    = errors.New("secret mismatch")
	ErrNameAlreadyExists = errors.New("a player with the same name already exists")
	ErrNameInvalid       = errors.New("invalid name")
)

func (d *Dispatcher) handleUnauthorized(socket *websocket.Conn, lobby *Lobby, data []byte) (err error) {
	str := strings.TrimSpace(string(data))
	spl := strings.Split(str, " ")

	switch spl[0] {
	case "JOIN":
		var name, password string
		// parse user:daniel pass:test token:ajcjvj
		for _, a := range spl[1:] {
			kv := strings.Split(a, ":")
			if len(kv) != 2 {
				continue
			}
			k, v := kv[0], kv[1]
			switch k {
			case "name":
				name = v
			case "pass", "password", "secret":
				password = v
			}
		}

		// check if a name or token was given
		// and if the name is valid
		if name == "" {
			return ErrArgs
		}
		if !IsNameValid(name) {
			return ErrNameInvalid
		}

		// check if the lobby is password protected
		if lobby.Secret != "" {
			if password != lobby.Secret {
				return ErrSecretMismatch
			}
		}

		// check if the name is already in the lobby
		if _, ok := lobby.Clients[strings.ToLower(name)]; ok {
			return ErrNameAlreadyExists
		}

		// create client and add to game
		client := NewClient(socket, name)
		lobby.Clients[strings.ToLower(name)] = client

		return d.Call(JoinEvent, socket, client, lobby, nil)
	}
	return
}

func (d *Dispatcher) handleMessage(socket *websocket.Conn, lobby *Lobby, data []byte) (err error) {
	// check payload size
	if len(data) > 4096 {
		return ErrPayloadTooLarge
	}

	// find client and allow authorized routes or allow unauthorized routes if not found
	_, client, ok := d.gobby.BySocket(socket)
	if !ok {
		return d.handleUnauthorized(socket, lobby, data)
	}

	// decode message
	var msg *Message
	if err = json.Unmarshal(data, &msg); err != nil {
		Warnf(socket, "cannot decode JSON message")
		return
	}
	msg.Command = strings.TrimSpace(strings.ToUpper(msg.Command))

	// check if the message is a reply
	if msg.To != "" {
		d.handleReply(socket, lobby, client, msg)
		return
	}

	// call event to listeners
	if err = d.Call(ReceiveEvent, socket, client, lobby, msg); err != nil {
		Warnf(socket, "cannot call receive event: %v", err)
		return
	}

	// run handler
	h, ok := d.gobby.Handlers[msg.Command]
	if !ok {
		return ErrHandlerNotFound
	}
	return h.Execute(socket, lobby, client, msg)
}

func (d *Dispatcher) handleReply(socket *websocket.Conn, lobby *Lobby, client *Client, msg *Message) {
	replyMessageHooksMu.RLock()
	reply, ok := replyMessageHooks[msg.To]
	replyMessageHooksMu.RUnlock()

	if !ok {
		Warnf(client, "tried to send a reply to an unknown message")
		return
	}

	// send reply and call event
	reply <- msg

	if err := d.Call(ReceiveReplyEvent, socket, client, lobby, msg); err != nil {
		Warnf(socket, "cannot call reply event: %v", err)
	}
	return
}

func (d *Dispatcher) handleClose(socket *websocket.Conn) {
	// check if socket was connected to a lobby
	lobby, client, ok := d.gobby.BySocket(socket)
	if !ok {
		// client wasn't important anyway
		return
	}
	d.gobby.RemoveClient(lobby, client)
}
