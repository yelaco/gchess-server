package game

import (
	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/pkg/utils"
)

type player struct {
	conn *websocket.Conn
	id   string
}

func GeneratePlayerIds() [2]string {
	return [2]string{utils.GenerateUUID(), utils.GenerateUUID()}
}
