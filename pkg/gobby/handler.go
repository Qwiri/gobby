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
	Handler func(*Handle) error
}

type Handle struct {
	Lobby   *Lobby
	Client  *Client
	Message *Message
	Args    validate.Result
}

func (h *Handle) String(name string) string {
	return h.Args.String(name)
}

func (h *Handle) Number(name string) int64 {
	return h.Args.Number(name)
}

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

func (h *Handler) Execute(socket *websocket.Conn, lobby *Lobby, client *Client, msg *Message) (err error) {
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
	// Check validation
	var res validate.Result
	if len(h.Validation) > 0 {
		if res, err = h.Validation.Check(msg.Args...); err != nil {
			return
		}
	}
	return h.Handler(&Handle{
		Lobby:   lobby,
		Client:  client,
		Message: msg,
		Args:    res,
	})
}
