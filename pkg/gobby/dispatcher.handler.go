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
	ErrCancelled         = errors.New("event cancelled")
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

		client := NewClient(socket, name)
		event := &Join{
			Client: client,
			Lobby:  lobby,
		}
		d.call(joinType, event)

		if event.cancelled {
			return ErrCancelled
		}

		lobby.Clients[strings.ToLower(name)] = client

		// send JOINED message to client to let client know the join was successful
		_ = NewBasicMessage("JOINED", name).Send(socket)
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

	// send raw message event
	d.call(messageReceiveRawType, &MessageReceiveRaw{
		Sender: client,
		Lobby:  lobby,
		Data:   data,
	})

	// decode message
	var msg *Message
	if err = json.Unmarshal(data, &msg); err != nil {
		Warnf(socket, "cannot decode JSON message")
		return
	}
	msg.Command = strings.TrimSpace(strings.ToLower(msg.Command))

	// check if the message is a reply
	if msg.To != "" {
		d.handleReply(lobby, client, msg)
		return
	}

	// run handler
	h, ok := d.gobby.Handlers[msg.Command]
	if !ok {
		return ErrHandlerNotFound
	}

	// call event to listeners
	d.call(messageReceiveType, &MessageReceive{
		Sender:  client,
		Lobby:   lobby,
		Message: msg,
		Handler: h,
	})

	return h.Execute(socket, lobby, client, msg)
}

func (d *Dispatcher) handleReply(lobby *Lobby, client *Client, msg *Message) {
	replyMessageHooksMu.RLock()
	reply, ok := replyMessageHooks[msg.To]
	replyMessageHooksMu.RUnlock()
	if !ok {
		Warnf(client, "tried to send a reply to an unknown message")
		return
	}
	reply <- msg
	d.call(messageReceiveReplyType, &MessageReceiveReply{
		Sender:  client,
		Lobby:   lobby,
		Message: msg,
	})
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
