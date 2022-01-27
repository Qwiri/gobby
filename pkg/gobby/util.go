package gobby

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gofiber/websocket/v2"
	"math/rand"
	"regexp"
	"strings"
)

var Version = "1.0.0"

func IsListener(i interface{}) bool {
	switch i.(type) {
	case BasicEvent, MessageEvent, AsyncMessageEvent, AsyncBasicEvent:
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

const CharSet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz123456789"

func generateRandomString(length int) string {
	var bob strings.Builder
	for i := 0; i < length; i++ {
		bob.WriteRune(rune(CharSet[rand.Intn(len(CharSet))]))
	}
	return bob.String()
}

///

var usernameExpr = regexp.MustCompile(`^[A-Za-z0-9_\-]{1,16}$`)

func IsNameValid(username string) bool {
	return usernameExpr.MatchString(username)
}
