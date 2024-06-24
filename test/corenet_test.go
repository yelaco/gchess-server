package test

import (
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/internal/corenet"
)

func TestWebSocketServer(t *testing.T) {
	wsServer := corenet.NewWebSocketServer()
	wsServer.SetMessageHandler(handleWebSocketMessage)
	wsServer.Start()
}

func handleWebSocketMessage(conn *websocket.Conn, message *corenet.Message) {
	fmt.Println(message)
	switch message.Action {
	case "matching":
		playerId, ok := message.Data["playerId"].(string)
		if ok {
			fmt.Println(playerId)
		}
	case "move":
		playerId, ok := message.Data["playerId"].(string)
		if !ok {
			return
		}
		sessionId, ok := message.Data["sessionId"].(string)
		if !ok {
			return
		}
		move := message.Data["move"].(string)
		if !ok {
			return
		}
		fmt.Println(playerId, sessionId, move)
	default:

	}
}
