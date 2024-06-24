package session

import (
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/internal/game"
	"github.com/yelaco/robinhood-chess/pkg/logging"
	"go.uber.org/zap"
)

type GameSession struct {
	Players map[string]*Player
	Game    *game.Game
}

var gameSessions = make(map[string]*GameSession)
var mu sync.Mutex

func InitSession(sessionID string, player1, player2 *Player) {
	playersMap := map[string]*Player{
		player1.ID: player1,
		player2.ID: player2,
	}
	gameSessions[sessionID] = &GameSession{
		Players: playersMap,
		Game:    game.InitGame([2]string{player1.ID, player2.ID}),
	}
}

func StartGame(session *GameSession) {
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
			type errorResponse struct {
				Error string `json:"error"`
			}
			logging.Warn("invalid move",
				zap.String("session_id", sessionID),
				zap.String("player_id", playerID),
				zap.String("move", move),
				zap.String("error", err.Error()),
			)
			session.Players[playerID].Conn.WriteJSON(errorResponse{
				Error: "invalid move: " + err.Error(),
			})
			return
		}

		logging.Info("valid move",
			zap.String("session_id", sessionID),
			zap.String("player_id", playerID),
			zap.String("move", move),
		)
		// notify players about the new board state
		for _, player := range session.Players {
			player.Conn.WriteJSON(session.Game.GetBoard())
			if session.Game.IsOver() {
				player.Conn.Close()
				delete(gameSessions, sessionID)
			}
		}
	}
	mu.Unlock()
}
