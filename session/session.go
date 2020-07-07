package session

import (
	"github.com/gofiber/session"
	"sync"
	"time"
)

var once sync.Once
var _session *session.Session

func Session() *session.Session {
	once.Do(func() {
		_session = session.New(session.Config{Expiration: time.Hour * 24})
	})
	return _session
}
