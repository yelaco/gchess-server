package matcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/yelaco/go-chess-server/pkg/config"
	"github.com/yelaco/go-chess-server/pkg/logging"
	"github.com/yelaco/go-chess-server/pkg/session"
	"github.com/yelaco/go-chess-server/pkg/utils"
	"go.uber.org/zap"
)

/*
A Matcher handles matchmaking logic and forwards the player connection to session manager
*/
type Matcher struct {
	Queue      []*session.Player
	SessionMap map[string]string
	ConnMap    map[string]string
	mu         sync.Mutex
}

type matchResponse struct {
	Type        string              `json:"type"`
	SessionID   string              `json:"session_id"`
	GameState   gameStateResponse   `json:"game_state"`
	PlayerState session.PlayerState `json:"player_state"`
}

type gameStateResponse struct {
	Status      string `json:"status,omitempty"`
	BoardFen    string `json:"board_fen,omitempty"`
	IsWhiteTurn bool   `json:"is_white_turn,omitempty"`
}

type timeoutResponpse struct {
	Type    string `json:"type"`
	Message string `json:"Message"`
}

/*
Return a Matcher with initialized fields
*/
func NewMatcher() *Matcher {
	return &Matcher{
		Queue:      []*session.Player{},
		SessionMap: map[string]string{},
		ConnMap:    map[string]string{},
		mu:         sync.Mutex{},
	}
}

/*
Enter players to the matching queue. Matcher also keeps track of connection ID
to ensure no user can enter queue multiple time at the same time.
After timeout, Matcher will cancel queueing of the corresponding player
if there aren't no matches available.
The player can also rejoin an unfinished match they left
*/
func (m *Matcher) EnterQueue(player *session.Player, connID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sessionID, exists := m.SessionMap[player.ID]
	if exists {
		m.rejoinMatch(sessionID, player)
		return
	}
	for _, pid := range m.ConnMap {
		if pid == player.ID {
			player.Conn.WriteJSON(struct {
				Type  string `json:"type"`
				Error string `json:"error"`
			}{
				Type:  "queueing",
				Error: "Already queued",
			})
			return
		}
	}
	m.Queue = append(m.Queue, player)
	m.ConnMap[connID] = player.ID
	go m.leaveQueueIfTimeout(player, connID)
	go m.findMatch()
}

/*
Matcher pushes player out of the matching queue after a timeout if there aren't no matches available.
*/
func (m *Matcher) leaveQueueIfTimeout(player *session.Player, connID string) {
	time.Sleep(config.MatchingTimeout)
	if player == nil {
		return
	}
	if _, ok := m.SessionMap[player.ID]; !ok {
		player.Conn.WriteJSON(timeoutResponpse{
			Type:    "timeout",
			Message: "Canceled matching due to timeout",
		})
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.ConnMap, connID)
	for i, p := range m.Queue {
		if p.ID == player.ID || p == player {
			m.Queue = append(m.Queue[:i], m.Queue[i+1:]...)
			return
		}
	}
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

		notifyMatchingResult(sessionID, player1)
		notifyMatchingResult(sessionID, player2)
	}
}

func (m *Matcher) rejoinMatch(sessionID string, player *session.Player) {
	if err := session.PlayerJoin(sessionID, player); err != nil {
		player.Conn.WriteJSON(struct {
			Type  string `json:"type"`
			Error string `json:"error"`
		}{
			Type:  "error",
			Error: "Coulnd't join match: " + err.Error(),
		})
		return
	}
	notifyMatchingResult(sessionID, player)
}

func notifyMatchingResult(sessionID string, player *session.Player) {
	gameState, err := session.GetGameState(sessionID)
	if err != nil {
		player.Conn.WriteJSON(struct {
			Type  string `json:"type"`
			Error string `json:"error"`
		}{
			Type:  "error",
			Error: "Coulnd't join match: " + err.Error(),
		})
		return
	}

	playerState, err := session.GetPlayerState(sessionID, player.ID)
	if err != nil {
		player.Conn.WriteJSON(struct {
			Type  string `json:"type"`
			Error string `json:"error"`
		}{
			Type:  "error",
			Error: "Coulnd't join match: " + err.Error(),
		})
	}

	player.Conn.WriteJSON(matchResponse{
		Type:      "matched",
		SessionID: sessionID,
		GameState: gameStateResponse{
			Status:      gameState.Status,
			BoardFen:    utils.BoardToFen(gameState.Board),
			IsWhiteTurn: gameState.IsWhiteTurn,
		},
		PlayerState: playerState,
	})
}

/*
Check if player is in a session
*/
func (m *Matcher) SessionExists(playerID string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sessionID, exists := m.SessionMap[playerID]
	if exists {
		return sessionID, true
	} else {
		return "", false
	}
}

/*
Remove the session after it terminated
*/
func (m *Matcher) RemoveSession(player1, player2 string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.SessionMap, player1)
	delete(m.SessionMap, player2)
}
