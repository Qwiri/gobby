package gobby

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gofiber/websocket/v2"
)

var Version = "1.0.0"

func IsListener(i interface{}) bool {
	switch i.(type) {
	case ByteEvent, MessagedEvent:
		return true
	}
	return false
}

func prefix(i interface{}) string {
	switch t := i.(type) {
	case *Client:
		return fmt.Sprintf("[client:%s]", t.Name)
	case *Lobby:
		return fmt.Sprintf("[lobby:%s]", t.ID)
	case *websocket.Conn:
		return fmt.Sprintf("{socket@%s}", t.RemoteAddr())
	case fmt.Stringer:
		return fmt.Sprintf("{str:%s}", t.String())
	}
	return "[?!]"
}

func Infof(i interface{}, msg string, args ...interface{}) {
	log.Infof(fmt.Sprintf("%s %s", prefix(i), msg), args...)
}

func Warnf(i interface{}, msg string, args ...interface{}) {
	log.Warnf(fmt.Sprintf("%s %s", prefix(i), msg), args...)
}
