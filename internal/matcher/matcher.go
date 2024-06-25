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
	Queue      []*session.Player
	SessionMap map[string]string
	mu         sync.Mutex
}

type matchResponse struct {
	Type        string       `json:"type"`
	SessionID   string       `json:"sessionId"`
	IsWhite     bool         `json:"isWhite"`
	Board       [8][8]string `json:"board"`
	IsWhiteTurn bool         `json:"isWhiteTurn"`
}

func NewMatcher() *Matcher {
	return &Matcher{
		Queue:      []*session.Player{},
		SessionMap: map[string]string{},
		mu:         sync.Mutex{},
	}
}

func (m *Matcher) EnterQueue(player *session.Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sessionID, exists := m.SessionMap[player.ID]
	if exists {
		m.RejoinMatch(sessionID, player)
		return
	}
	m.Queue = append(m.Queue, player)
	go m.findMatch()
}

func generateSessionId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (m *Matcher) findMatch() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.Queue) >= 2 {
		player1 := m.Queue[0]
		player2 := m.Queue[1]
		m.Queue = m.Queue[2:]

		sessionID := generateSessionId()
		session.InitSession(sessionID, player1, player2)
		m.SessionMap[player1.ID] = sessionID
		m.SessionMap[player2.ID] = sessionID

		logging.Info("init match",
			zap.String("player_1", player1.ID),
			zap.String("player_2", player2.ID),
		)

		// Notify players the new session
		player1.Conn.WriteJSON(matchResponse{
			Type:      "matched",
			SessionID: sessionID,
			IsWhite:   true,
		})
		player2.Conn.WriteJSON(matchResponse{
			Type:      "matched",
			SessionID: sessionID,
			IsWhite:   false,
		})
	}
}

func (m *Matcher) RejoinMatch(sessionID string, player *session.Player) {
	isWhiteSide, err := session.GetPlayerSide(sessionID, player.ID)
	if err != nil {
		player.Conn.WriteJSON(struct {
			Type  string `json:"type"`
			Error string `json:"error"`
		}{
			Type:  "error",
			Error: "Coulnd't rejoin match: " + err.Error(),
		})
		return
	}

	board, isWhiteTurn, err := session.GetGameState(sessionID)
	if err != nil {
		player.Conn.WriteJSON(struct {
			Type  string `json:"type"`
			Error string `json:"error"`
		}{
			Type:  "error",
			Error: "Coulnd't rejoin match: " + err.Error(),
		})
		return
	}

	session.PlayerJoin(sessionID, player)

	player.Conn.WriteJSON(matchResponse{
		Type:        "matched",
		SessionID:   sessionID,
		IsWhite:     isWhiteSide,
		Board:       board,
		IsWhiteTurn: isWhiteTurn,
	})
}

func (m *Matcher) RemoveSession(player1, player2 string) {
	delete(m.SessionMap, player1)
	delete(m.SessionMap, player2)
}
