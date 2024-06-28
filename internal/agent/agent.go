package agent

import (
	"github.com/gorilla/websocket"
	"github.com/yelaco/go-chess-server/internal/corenet"
	"github.com/yelaco/go-chess-server/internal/database"
	"github.com/yelaco/go-chess-server/internal/matcher"
	"github.com/yelaco/go-chess-server/internal/session"
	"github.com/yelaco/go-chess-server/pkg/logging"
	"go.uber.org/zap"
)

type Agent struct {
	wsServer *corenet.WebSocketServer
	matcher  *matcher.Matcher
}

func NewAgent() *Agent {
	a := &Agent{
		wsServer: corenet.NewWebSocketServer(),
		matcher:  matcher.NewMatcher(),
	}
	a.wsServer.SetMessageHandler(a.handleWebSocketMessage)
	a.wsServer.SetConnCloseGameHandler(a.playerDisconnectHandler)
	session.SetGameOverHandler(a.handleSessionGameOver)

	return a
}

func (a *Agent) StartGameServer() error {
	err := a.wsServer.Start()
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) handleSessionGameOver(s *session.GameSession, sessionID string) {
	playerIDs := make([]string, 0, 2)
	for id, player := range s.Players {
		playerIDs = append(playerIDs, id)
		player.Conn.WriteJSON(struct {
			Type string            `json:"type"`
			Data map[string]string `json:"data"`
		}{
			Type: "endgame",
			Data: map[string]string{
				"game_state": s.Game.GetStatus(),
			},
		})
		player.Conn.Close()
	}
	gameMoves := s.Game.GetAllMoves()
	if _, err := database.InsertSession(sessionID, playerIDs[0], playerIDs[1], gameMoves); err != nil {
		logging.Error("coulnd't save game", zap.Error(err))
	}
	session.CloseSession(sessionID)
	a.matcher.RemoveSession(playerIDs[0], playerIDs[1])
}

func (a *Agent) playerDisconnectHandler(playerID string) {
	sessionID, exists := a.matcher.SessionExists(playerID)
	if !exists {
		return
	}

	err := session.PlayerLeave(sessionID, playerID)
	if err != nil {
		logging.Warn("player disconnected error",
			zap.String("player_id", playerID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
	}

	logging.Info("player disconnected",
		zap.String("player_id", playerID),
		zap.String("session_id", sessionID),
	)
}

func (a *Agent) handleWebSocketMessage(conn *websocket.Conn, message *corenet.Message, connID *string) {
	type errorResponse struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}
	switch message.Action {
	case "matching":
		playerID, ok := message.Data["player_id"].(string)
		if ok {
			*connID = playerID
			logging.Info("attempt matchmaking",
				zap.String("status", "queued"),
				zap.String("player_id", playerID),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			a.matcher.EnterQueue(&session.Player{
				Conn: conn,
				ID:   playerID,
			})
		} else {
			logging.Info("attempt matchmaking",
				zap.String("status", "rejected"),
				zap.String("error", "insufficient data"),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			conn.WriteJSON(errorResponse{
				Type:  "error",
				Error: "insufficient data",
			})
		}
	case "move":
		playerID, playerOK := message.Data["player_id"].(string)
		sessionID, sessionOK := message.Data["session_id"].(string)
		move, moveOK := message.Data["move"].(string)
		if playerOK && sessionOK && moveOK {
			logging.Info("attempt making move",
				zap.String("status", "processing"),
				zap.String("player_id", playerID),
				zap.String("session_id", sessionID),
				zap.String("move", move),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			session.ProcessMove(sessionID, playerID, move)
		} else {
			logging.Info("attempt making move",
				zap.String("status", "rejected"),
				zap.String("error", "insufficient data"),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			conn.WriteJSON(errorResponse{
				Type:  "error",
				Error: "insufficient data",
			})
		}
	default:
	}
}
