package session

import "github.com/gorilla/websocket"

type Player struct {
	Conn *websocket.Conn
	ID   string `json:"id"`
}
