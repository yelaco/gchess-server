package agent

import (
	"github.com/gorilla/websocket"
	"github.com/yelaco/robinhood-chess/internal/corenet"
	"github.com/yelaco/robinhood-chess/internal/database"
	"github.com/yelaco/robinhood-chess/internal/matcher"
	"github.com/yelaco/robinhood-chess/internal/session"
	"github.com/yelaco/robinhood-chess/pkg/logging"
	"go.uber.org/zap"
)

type Agent struct {
	db       *database.DBConnection
	wsServer *corenet.WebSocketServer
	matcher  *matcher.Matcher
}

func NewAgent() *Agent {
	dbConn, err := database.Connect()
	if err != nil {
		logging.Fatal("NewAgent(): couldn't init database connection")
	}

	a := &Agent{
		db:       dbConn,
		wsServer: corenet.NewWebSocketServer(),
		matcher:  matcher.NewMatcher(),
	}
	a.wsServer.SetMessageHandler(a.handleWebSocketMessage)

	return a
}

func (a *Agent) StartGameServer() error {
	err := a.wsServer.Start()
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) Close() {
	a.db.Close()
}

func (a *Agent) handleWebSocketMessage(conn *websocket.Conn, message *corenet.Message) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	switch message.Action {
	case "matching":
		playerId, ok := message.Data["playerId"].(string)
		if ok {
			logging.Info("attempt matchmaking",
				zap.String("status", "queued"),
				zap.String("player_id", playerId),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			a.matcher.EnterQueue(&session.Player{
				Conn: conn,
				ID:   playerId,
			})
		} else {
			logging.Info("attempt matchmaking",
				zap.String("status", "rejected"),
				zap.String("error", "insufficient data"),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			conn.WriteJSON(errorResponse{
				Error: "insufficient data",
			})
		}
	case "move":
		playerId, playerOK := message.Data["playerId"].(string)
		sessionId, sessionOK := message.Data["sessionId"].(string)
		move, moveOK := message.Data["move"].(string)
		if playerOK && sessionOK && moveOK {
			logging.Info("attempt making move",
				zap.String("status", "processing"),
				zap.String("player_id", playerId),
				zap.String("session_id", sessionId),
				zap.String("move", move),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			session.ProcessMove(sessionId, playerId, move)
		} else {
			logging.Info("attempt making move",
				zap.String("status", "rejected"),
				zap.String("error", "insufficient data"),
				zap.String("remote_address", conn.RemoteAddr().String()),
			)
			conn.WriteJSON(errorResponse{
				Error: "insufficient data",
			})
		}
	default:

	}
}
