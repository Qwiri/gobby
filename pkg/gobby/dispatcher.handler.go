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
	ErrMessageFormat   = errors.New("message format error")

	ErrArgs              = errors.New("args mismatch")
	ErrSecretMismatch    = errors.New("secret mismatch")
	ErrNameAlreadyExists = errors.New("a player with the same name already exists")
	ErrNameInvalid       = errors.New("invalid name")
	ErrPassInvalid       = errors.New("invalid password")
	ErrCancelled         = errors.New("event cancelled")
)

func (d *Dispatcher) handleUnauthorized(socket *websocket.Conn, lobby *Lobby, msg *Message) (err error) {
	// unauthorized handlers
	switch msg.Command {
	case "JOIN":
		var (
			username string
			password string
			ok       bool
		)
		if len(msg.Args) <= 0 {
			return msg.ReplyError(socket, ErrArgs)
		}

		// parse name and check validity
		if username, ok = msg.Args[0].(string); !ok {
			return msg.ReplyError(socket, ErrNameInvalid)
		}
		if !IsNameValid(username) {
			return msg.ReplyError(socket, ErrNameInvalid)
		}

		// check if the lobby is password protected
		if lobby.Secret != "" {
			if len(msg.Args) > 1 {
				if password, ok = msg.Args[1].(string); !ok {
					return msg.ReplyError(socket, ErrPassInvalid)
				}
			}
			if password != lobby.Secret {
				return msg.ReplyError(socket, ErrSecretMismatch)
			}
		}

		// check if the name is already in the lobby
		if _, ok = lobby.Clients[strings.ToLower(username)]; ok {
			return msg.ReplyError(socket, ErrNameAlreadyExists)
		}

		client := NewClient(socket, username)

		// call join event which can be cancelled
		event := &Join{
			Client:  client,
			Lobby:   lobby,
			Message: msg,
		}
		d.call(joinType, event)

		if event.cancelled {
			return msg.ReplyError(socket, ErrCancelled)
		}

		lobby.Clients[strings.ToLower(username)] = client

		// send JOINED message to client to let client know the join was successful
		return msg.ReplyBasic(socket, "JOINED", username)
	}
	return nil
}

func (d *Dispatcher) handleAuthorized(socket *websocket.Conn, lobby *Lobby, msg *Message, client *Client) (err error) {
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

func (d *Dispatcher) handleMessage(socket *websocket.Conn, lobby *Lobby, data []byte) (err error) {
	// check payload size
	if len(data) > 4096 {
		return ErrPayloadTooLarge
	}

	// decode message
	var msg *Message
	if err = json.Unmarshal(data, &msg); err != nil {
		Warnf(socket, "cannot decode JSON message")
		return ErrMessageFormat
	}
	msg.Command = strings.TrimSpace(strings.ToLower(msg.Command))

	// find client and allow authorized routes or allow unauthorized routes if not found
	_, client, ok := d.gobby.BySocket(socket)
	if ok {
		return d.handleAuthorized(socket, lobby, msg, client)
	} else {
		return d.handleUnauthorized(socket, lobby, msg)
	}
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
