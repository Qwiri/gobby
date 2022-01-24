package gobby

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
)

func (d *Dispatcher) handleUnauthorized(socket *websocket.Conn, data []byte) {
	// TODO: Implement this
}

func (d *Dispatcher) handleMessage(socket *websocket.Conn, data []byte) {
	// find client and allow authorized routes or allow unauthorized routes if not found
	lobby, client, ok := d.gobby.BySocket(socket)
	if !ok {
		d.handleUnauthorized(socket, data)
		return
	}

	// decode message
	var msg *Message
	if err := json.Unmarshal(data, &msg); err != nil {
		Warnf(socket, "cannot decode JSON message")
		return
	}

	// check if the message is a reply
	if msg.To != "" {
		d.handleReply(socket, lobby, client, msg)
		return
	}

	// call event to listeners
	if err := d.Call(ReceiveEvent, socket, client, lobby, msg); err != nil {
		Warnf(socket, "cannot call receive event: %v", err)
		return
	}
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
