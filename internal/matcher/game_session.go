package matcher

import (
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/internal/game"
	"github.com/yelaco/robinhood-chess/internal/logging"
)

type GameSession struct {
	Players []*Player
	Game    *game.Game
}

var gameSessions = make(map[string]*GameSession)
var mu sync.Mutex

func StartGame(session *GameSession) {
	session.Game = game.InitGame()
	for _, player := range session.Players {
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"start"}`))
		if err != nil {
			log.Println("Error sending start message:", err)
		}
	}
}

func ProcessMove(sessionID, playerID, move string) {
	mu.Lock()
	session, exists := gameSessions[sessionID]
	if exists {
		pos := strings.Split(move, "-")
		err := session.Game.MakeMove(playerID, pos[0], pos[1])
		if err != nil {
			logging.Error("Can't make move")
			return
		}

		// notify players about the new board state
	}
	mu.Unlock()
}
