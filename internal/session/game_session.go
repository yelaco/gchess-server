package session

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yelaco/go-chess-server/internal/game"
	"github.com/yelaco/go-chess-server/pkg/logging"
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

type SessionResponse struct {
	Type      string    `json:"type"`
	GameState GameState `json:"game_state"`
}

type PlayerState struct {
	IsWhiteSide bool `json:"is_white_side"`
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
	mu.Lock()
	defer mu.Unlock()
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

func PlayerInSession(sessionID string, player *Player) bool {
	mu.Lock()
	defer mu.Unlock()
	session, exists := gameSessions[sessionID]
	if exists {
		if p, ok := session.Players[player.ID]; ok {
			return p == player
		}
		return false
	}
	return false
}

func PlayerJoin(sessionID string, player *Player) error {
	mu.Lock()
	defer mu.Unlock()
	session, exists := gameSessions[sessionID]
	if exists {
		if p, ok := session.Players[player.ID]; ok {
			if p != nil {
				return errors.New("player still in session")
			}
			session.Players[player.ID] = player
			return nil
		}
		return errors.New("player id not in the session")
	}
	return errors.New("invalid session id")
}

func PlayerLeave(sessionID, playerID string) error {
	mu.Lock()
	defer mu.Unlock()
	session, exists := gameSessions[sessionID]
	if exists {
		session.Players[playerID] = nil
		return nil
	}
	return errors.New("invalid session id")
}

func ProcessMove(sessionID, playerID, move string) {
	mu.Lock()

	type errorResponse struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}

	session, exists := gameSessions[sessionID]
	if exists {
		pos, parseErr := game.ParseMove(move)
		if parseErr != nil {
			logging.Warn("invalid move",
				zap.String("session_id", sessionID),
				zap.String("player_id", playerID),
				zap.String("move", move),
				zap.String("error", parseErr.Error()),
			)
			if err := session.Players[playerID].Conn.WriteJSON(errorResponse{
				Type:  "error",
				Error: "invalid move: " + parseErr.Error(),
			}); err != nil {
				logging.Info("ws write", zap.Error(err))
			}
			mu.Unlock()
			return
		}

		err := session.Game.MakeMove(playerID, pos[0], pos[1])
		if err != nil {
			logging.Warn("invalid move",
				zap.String("session_id", sessionID),
				zap.String("player_id", playerID),
				zap.String("move", move),
				zap.String("error", err.Error()),
			)
			if err := session.Players[playerID].Conn.WriteJSON(errorResponse{
				Type:  "error",
				Error: "invalid move: " + err.Error(),
			}); err != nil {
				logging.Info("ws write", zap.Error(err))
			}
			mu.Unlock()
			return
		}

		logging.Info("valid move",
			zap.String("session_id", sessionID),
			zap.String("player_id", playerID),
			zap.String("move", move),
		)

		mu.Unlock()

		// notify players about the new board state
		for _, player := range session.Players {
			gameState, err := GetGameState(sessionID)
			if err != nil {
				logging.Error("invalid session id for game state")
				if err := player.Conn.WriteJSON(errorResponse{
					Type:  "error",
					Error: "coulnd't retrieve game state",
				}); err != nil {
					logging.Info("ws write", zap.Error(err))
				}
				return
			}

			if player == nil {
				continue
			}

			if err := player.Conn.WriteJSON(SessionResponse{
				Type:      "session",
				GameState: gameState,
			}); err != nil {
				logging.Error("couldn't notify player ", zap.String("player_id", playerID))
			}
		}

		if session.Game.IsOver() {
			gameOverHandler(session, sessionID)
		}
	}
}
