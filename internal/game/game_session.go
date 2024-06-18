package game

import (
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/internal/logging"
)

type GameSession struct {
	players []player
	game    *Game
}

var gameSessions = make(map[string]*GameSession)
var mu sync.Mutex

func StartGame(session *GameSession) {
	session.game = InitGame()
	for _, player := range session.players {
		err := player.conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"start"}`))
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
		err := session.game.MakeMove(playerID, pos[0], pos[1])
		if err != nil {
			logging.Error("Can't make move")
			return
		}

		// notify players about the new board state
	}
	mu.Unlock()
}

func InitGame() *Game {
	g := &Game{
		playerIds:   GeneratePlayerIds(),
		board:       initBoard(),
		isWhiteTurn: true,
		status:      active,
	}
	g.kingSpots[0] = g.board.boxes[4][0]
	g.kingSpots[1] = g.board.boxes[4][7]
	return g
}
