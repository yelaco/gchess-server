package session

import (
	"errors"
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

type GameState struct {
	Status      string       `json:"status"`
	Board       [8][8]string `json:"board"`
	IsWhiteTurn bool         `json:"is_white"`
}

type PlayerState struct {
	IsWhiteSide bool
}

var gameSessions = make(map[string]*GameSession)
var mu sync.RWMutex
var gameOverHandler = func(session *GameSession, sessionID string) {
	CloseSession(sessionID)
	for _, player := range session.Players {
		player.Conn.Close()
	}
}

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

func CloseSession(sessionID string) {
	delete(gameSessions, sessionID)
}

func SetGameOverHandler(govHandler func(*GameSession, string)) {
	gameOverHandler = govHandler
}

func StartGame(session *GameSession) {
	for _, player := range session.Players {
		err := player.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"start"}`))
		if err != nil {
			log.Println("Error sending start message:", err)
		}
	}
}

func GetGameState(sessionID string) (GameState, error) {
	mu.RLock()
	defer mu.RUnlock()
	session, exists := gameSessions[sessionID]
	if exists {
		return GameState{
			Status:      session.Game.GetStatus(),
			Board:       session.Game.GetBoard(),
			IsWhiteTurn: session.Game.GetCurrentTurn(),
		}, nil
	}
	return GameState{}, errors.New("invalid session id")
}

func GetPlayerState(sessionID, playerID string) (PlayerState, error) {
	mu.RLock()
	defer mu.RUnlock()
	session, exists := gameSessions[sessionID]
	if exists {
		isWhiteSide, err := session.Game.GetPlayerSide(playerID)
		if err != nil {
			return PlayerState{}, errors.New("invalid player id")
		}
		return PlayerState{
			IsWhiteSide: isWhiteSide,
		}, nil
	}
	return PlayerState{}, errors.New("invalid session id")
}

func PlayerJoin(sessionID string, player *Player) error {
	session, exists := gameSessions[sessionID]
	if exists {
		if _, ok := session.Players[player.ID]; ok {
			session.Players[player.ID] = player
			return nil
		}
		return errors.New("player id not in the session")
	}
	return errors.New("invalid session id")
}

func PlayerLeave(sessionID, playerID string) error {
	session, exists := gameSessions[sessionID]
	if exists {
		session.Players[playerID] = nil
		return nil
	}
	return errors.New("invalid session id")
}

func ProcessMove(sessionID, playerID, move string) {
	mu.Lock()
	defer mu.Unlock()
	session, exists := gameSessions[sessionID]
	if exists {
		pos := strings.Split(move, "-")
		err := session.Game.MakeMove(playerID, pos[0], pos[1])
		if err != nil {
			type errorResponse struct {
				Type  string `json:"type"`
				Error string `json:"error"`
			}
			logging.Warn("invalid move",
				zap.String("session_id", sessionID),
				zap.String("player_id", playerID),
				zap.String("move", move),
				zap.String("error", err.Error()),
			)
			session.Players[playerID].Conn.WriteJSON(errorResponse{
				Type:  "error",
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
		}

		if session.Game.IsOver() {
			gameOverHandler(session, sessionID)
		}
	}
}
