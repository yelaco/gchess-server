package matcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/yelaco/robinhood-chess/internal/session"
	"github.com/yelaco/robinhood-chess/pkg/logging"
	"go.uber.org/zap"
)

type Matcher struct {
	Queue []*session.Player
	mu    sync.Mutex
}

func NewMatcher() *Matcher {
	return &Matcher{
		Queue: []*session.Player{},
		mu:    sync.Mutex{},
	}
}

func (m *Matcher) EnterQueue(player *session.Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Queue = append(m.Queue, player)
	go m.findMatch()
}

func generateSessionId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (m *Matcher) findMatch() {
	m.mu.Lock()
	defer m.mu.Unlock()

	type gameResponse struct {
		SessionID string `json:"sessionId"`
		IsWhite   bool   `json:"isWhite"`
	}

	if len(m.Queue) >= 2 {
		player1 := m.Queue[0]
		player2 := m.Queue[1]
		m.Queue = m.Queue[2:]

		sessionID := generateSessionId()
		session.InitSession(sessionID, player1, player2)

		logging.Info("init match",
			zap.String("player_1", player1.ID),
			zap.String("player_2", player2.ID),
		)

		// Notify players the new session
		player1.Conn.WriteJSON(gameResponse{
			SessionID: sessionID,
			IsWhite:   true,
		})
		player2.Conn.WriteJSON(gameResponse{
			SessionID: sessionID,
			IsWhite:   false,
		})
	}
}

// func (m *Matcher) JoinSession(conn *websocket.Conn, sessionID string) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	session, exists := gameSessions[sessionID]
// 	if !exists {
// 	} else {
// 		session.Players = append(session.Players, conn)
// 		if len(session.Players) == 2 {
// 			StartGame(session)
// 		}
// 	}
// }
