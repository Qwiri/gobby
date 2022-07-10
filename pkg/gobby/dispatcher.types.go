package gobby

const (
	lobbyCreateType eventType = iota //
	joinType                         //
	leaveType
	messageReceiveRawType   //
	messageReceiveType      //
	messageReceiveReplyType //
)

func (d *Dispatcher) getEventHandler(i interface{}) eventHandler {
	switch t := i.(type) {
	case func(event *LobbyCreate):
		return lobbyCreateHandler(t)
	case func(event *Join):
		return joinHandler(t)
	case func(event *Leave):
		return leaveHandler(t)
	case func(event *MessageReceiveRaw):
		return receiveRawHandler(t)
	case func(event *MessageReceive):
		return receiveHandler(t)
	case func(event *MessageReceiveReply):
		return receiveReplyHandler(t)
	}
	return nil
}

type lobbyCreateHandler func(event *LobbyCreate)

func (l lobbyCreateHandler) Type() eventType {
	return lobbyCreateType
}
func (l lobbyCreateHandler) Handle(i interface{}) {
	if v, o := i.(*LobbyCreate); o {
		l(v)
	}
}

///

type joinHandler func(event *Join)

func (j joinHandler) Type() eventType {
	return joinType
}
func (j joinHandler) Handle(i interface{}) {
	if v, o := i.(*Join); o {
		j(v)
	}
}

///

type leaveHandler func(event *Leave)

func (l leaveHandler) Type() eventType {
	return leaveType
}
func (l leaveHandler) Handle(i interface{}) {
	if v, o := i.(*Leave); o {
		l(v)
	}
}

///

type receiveRawHandler func(event *MessageReceiveRaw)

func (r receiveRawHandler) Type() eventType {
	return messageReceiveRawType
}
func (r receiveRawHandler) Handle(i interface{}) {
	if v, o := i.(*MessageReceiveRaw); o {
		r(v)
	}
}

///

type receiveReplyHandler func(event *MessageReceiveReply)

func (r receiveReplyHandler) Type() eventType {
	return messageReceiveReplyType
}
func (r receiveReplyHandler) Handle(i interface{}) {
	if v, o := i.(*MessageReceiveReply); o {
		r(v)
	}
}

///

type receiveHandler func(event *MessageReceive)

func (r receiveHandler) Type() eventType {
	return messageReceiveType
}
func (r receiveHandler) Handle(i interface{}) {
	if v, o := i.(*MessageReceive); o {
		r(v)
	}
}

// message receive

type LobbyCreate struct {
	Lobby *Lobby
	Addr  string
}

type Join struct {
	Client    *Client
	Lobby     *Lobby
	cancelled bool
}

func (j *Join) Cancel() {
	j.cancelled = true
}

type Leave struct {
	Client *Client
	Lobby  *Lobby
}

type MessageReceive struct {
	Sender  *Client
	Lobby   *Lobby
	Message *Message
	Handler *Handler
}

type MessageReceiveReply struct {
	Sender  *Client
	Lobby   *Lobby
	Message *Message
}

type MessageReceiveRaw struct {
	Sender *Client
	Lobby  *Lobby
	Data   []byte
}
