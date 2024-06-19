package matcher

import (
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn
	ID       string
	TimeLeft time.Duration
}
