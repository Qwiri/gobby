package gobby

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type LobbyID string

type Client struct {
	Name   string
	Socket *websocket.Conn
	Role   Role
	Meta   interface{}
}

func NewClient(socket *websocket.Conn, name string) *Client {
	return &Client{
		Name:   name,
		Socket: socket,
		Role:   0,
	}
}

///

type Lobby struct {
	ID      LobbyID
	State   State
	Clients map[string]*Client
	Secret  string
	Meta    interface{}
	lock    sync.RWMutex
}

func NewLobby(id LobbyID) *Lobby {
	return &Lobby{
		ID:      id,
		State:   0,
		Clients: make(map[string]*Client),
	}
}

// ChangeState sets a newState for a lobby and broadcasts the state change to all clients with
// the old state[0] and the new[1] state als arguments
func (l *Lobby) ChangeState(newState State) {
	oldState := l.State
	l.State = newState
	l.BroadcastForce(NewBasicMessage("STATE_CHANGE", oldState, newState))
}

type Router struct {
	g   *Gobby
	app *fiber.App
}
