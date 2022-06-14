package gobby

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Qwiri/gobby/pkg/util"
	"github.com/apex/log"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"sync"
	"time"
)

type MessageID string

type Message struct {
	// ID contains the message ID
	ID MessageID `json:"id"`
	// To can be set if the message is a reply to another message
	To MessageID `json:"to"`
	// Command contains the specific handler
	Command string `json:"cmd"`
	// Args contains (optional) arguments
	Args    []interface{} `json:"args"`
	Respond bool          `json:"respond"`
}

var (
	replyMessageHooks   = make(map[MessageID]chan<- *Message)
	replyMessageHooksMu sync.RWMutex
)
var ErrorReplyTimeout = errors.New("reply timeout")

func (m *Message) Marshal() (res []byte) {
	var err error
	if res, err = json.Marshal(m); err != nil {
		log.WithField("message", m).Warn("cannot marshal message")
	}
	return
}

func (m *Message) SendTo(client *Client) (err error) {
	err = client.Socket.WriteMessage(websocket.TextMessage, m.Marshal())
	return
}

func (m *Message) Send(socket *websocket.Conn) (err error) {
	err = socket.WriteMessage(websocket.TextMessage, m.Marshal())
	return
}

func (m *Message) SendAndAwaitReply(socket *websocket.Conn, timeout time.Duration) (*Message, error) {
	// mark Message to require reply
	m.Respond = true

	// create a channel where we can receive replies
	cr := make(chan *Message, 1)
	defer close(cr)

	replyMessageHooksMu.Lock()
	replyMessageHooks[m.ID] = cr
	replyMessageHooksMu.Unlock()

	defer func() {
		replyMessageHooksMu.Lock()
		delete(replyMessageHooks, m.ID)
		replyMessageHooksMu.Unlock()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Info("created timeout")

	// send message to client
	if err := m.Send(socket); err != nil {
		return nil, err
	}

	log.Info("waiting for reply")

	select {
	case res := <-cr:
		return res, nil
	case <-ctx.Done():
		return nil, ErrorReplyTimeout
	}
}

func NewBasicMessage(cmd string, args ...interface{}) *Message {
	return &Message{
		// generate UUID v4
		ID:      MessageID(uuid.New().String()),
		Command: cmd,
		Args:    args,
	}
}
func NewBasicMessageWith[T any](cmd string, args ...T) *Message {
	return NewBasicMessage(cmd, util.ConvertArray(args)...)
}

func NewErrorMessage(err error) *Message {
	return NewBasicMessage("ERROR", err.Error())
}

func (m *Message) CreateReply(msg *Message) *Message {
	msg.To = m.ID
	return msg
}

func (m *Message) ReplyWith(client *Client, msg Message) error {
	msg.To = m.ID
	return msg.SendTo(client)
}
