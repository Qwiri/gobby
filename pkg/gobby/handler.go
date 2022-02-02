package gobby

import (
	"errors"
	"github.com/Qwiri/gobby/pkg/validate"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	// Roles which have access to the handler (Optional)
	Roles Role
	// States (optional)
	States State
	// Validation to check the args for
	Validation validate.Schemes
	// Handler should be either BasicHandler or MessagedHandler
	// Required
	Handler interface{}
}

type (
	BasicHandler    func(*websocket.Conn, *Lobby, *Client) error
	MessagedHandler func(*websocket.Conn, *Lobby, *Client, *Message) error
)

type (
	Role  uint8
	State uint8
)

const (
	RoleLeader Role = 1 << iota
)

var (
	ErrStateNotAllowed = errors.New("handler cannot be used in this state")
	ErrRoleNotAllowed  = errors.New("handler cannot be used with this role")
)

func (h *Handler) Execute(socket *websocket.Conn, lobby *Lobby, client *Client, msg *Message) error {
	// Check State
	if h.States != 0 {
		if lobby.State&h.States != lobby.State {
			return ErrStateNotAllowed
		}
	}
	// Check Role
	if h.Roles != 0 {
		if client.Role&h.Roles != client.Role {
			return ErrRoleNotAllowed
		}
	}
	switch t := h.Handler.(type) {
	case BasicHandler:
		return t(socket, lobby, client)
	case MessagedHandler:
		return t(socket, lobby, client, msg)
	default:
		panic("invalid handler type")
	}
}
