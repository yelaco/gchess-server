package matcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/yelaco/robinhood-chess/internal/game"
)

type Matcher struct {
	Queue    []*Player
	Sessions map[string]*GameSession
	mu       sync.Mutex
}

func NewMatcher() *Matcher {
	return &Matcher{
		Queue:    []*Player{},
		Sessions: make(map[string]*GameSession),
	}
}

func (m *Matcher) EnterQueue(player *Player) {
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

	if len(m.Queue) >= 2 {
		player1 := m.Queue[0]
		player2 := m.Queue[1]
		m.Queue = m.Queue[2:]

		sessionID := generateSessionId()
		m.Sessions[sessionID] = &GameSession{
			Players: []*Player{player1, player2},
			Game:    game.InitGame(),
		}
		fmt.Printf("Match found: %s vs %s\n", player1.ID, player2.ID)

		// Notify players and game server about the new session
	}
}
