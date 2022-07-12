package gobby

const (
	lobbyCreateType eventType = iota //
	preJoinType                      //
	postJoinType
	leaveType
	messageReceiveRawType   //
	messageReceiveType      //
	messageReceiveReplyType //
)

func (d *Dispatcher) getEventHandler(i interface{}) eventHandler {
	switch t := i.(type) {
	case func(event *LobbyCreateEvent):
		return lobbyCreateHandler(t)
	case func(event *PreJoinEvent):
		return preJoinHandler(t)
	case func(event *PostJoinEvent):
		return postJoinHandler(t)
	case func(event *LeaveEvent):
		return leaveHandler(t)
	case func(event *MessageReceiveRawEvent):
		return receiveRawHandler(t)
	case func(event *MessageReceiveEvent):
		return receiveHandler(t)
	case func(event *MessageReceiveReplyEvent):
		return receiveReplyHandler(t)
	}
	return nil
}

type lobbyCreateHandler func(event *LobbyCreateEvent)

func (l lobbyCreateHandler) Type() eventType {
	return lobbyCreateType
}
func (l lobbyCreateHandler) Handle(i interface{}) {
	if v, o := i.(*LobbyCreateEvent); o {
		l(v)
	}
}

///

type preJoinHandler func(event *PreJoinEvent)

func (j preJoinHandler) Type() eventType {
	return preJoinType
}
func (j preJoinHandler) Handle(i interface{}) {
	if v, o := i.(*PreJoinEvent); o {
		j(v)
	}
}

///

type postJoinHandler func(event *PostJoinEvent)

func (p postJoinHandler) Type() eventType {
	return postJoinType
}

func (p postJoinHandler) Handle(i interface{}) {
	if v, o := i.(*PostJoinEvent); o {
		p(v)
	}
}

///

type leaveHandler func(event *LeaveEvent)

func (l leaveHandler) Type() eventType {
	return leaveType
}
func (l leaveHandler) Handle(i interface{}) {
	if v, o := i.(*LeaveEvent); o {
		l(v)
	}
}

///

type receiveRawHandler func(event *MessageReceiveRawEvent)

func (r receiveRawHandler) Type() eventType {
	return messageReceiveRawType
}
func (r receiveRawHandler) Handle(i interface{}) {
	if v, o := i.(*MessageReceiveRawEvent); o {
		r(v)
	}
}

///

type receiveReplyHandler func(event *MessageReceiveReplyEvent)

func (r receiveReplyHandler) Type() eventType {
	return messageReceiveReplyType
}
func (r receiveReplyHandler) Handle(i interface{}) {
	if v, o := i.(*MessageReceiveReplyEvent); o {
		r(v)
	}
}

///

type receiveHandler func(event *MessageReceiveEvent)

func (r receiveHandler) Type() eventType {
	return messageReceiveType
}
func (r receiveHandler) Handle(i interface{}) {
	if v, o := i.(*MessageReceiveEvent); o {
		r(v)
	}
}

// message receive

type LobbyCreateEvent struct {
	Lobby *Lobby
	Addr  string
}

// PreJoinEvent is called when a player tries to JOIN a lobby
// if canceled, the player receives an error on join and the socket is closed
type PreJoinEvent struct {
	Client    *Client
	Lobby     *Lobby
	Message   *Message
	cancelled bool
}

func (j *PreJoinEvent) Cancel() {
	j.cancelled = true
}

// PostJoinEvent is called after a player successfully JOINed a lobby
type PostJoinEvent struct {
	Client  *Client
	Lobby   *Lobby
	Message *Message
}

type LeaveEvent struct {
	Client *Client
	Lobby  *Lobby
}

type MessageReceiveEvent struct {
	Sender  *Client
	Lobby   *Lobby
	Message *Message
	Handler *Handler
}

type MessageReceiveReplyEvent struct {
	Sender  *Client
	Lobby   *Lobby
	Message *Message
}

type MessageReceiveRawEvent struct {
	Sender *Client
	Lobby  *Lobby
	Data   []byte
}
