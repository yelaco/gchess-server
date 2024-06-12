package test

import (
	"testing"

	"github.com/minhquang053/robinhood-chess/game"
	"github.com/minhquang053/robinhood-chess/internal/logging"
)

func TestGame(t *testing.T) {
	igame := game.InitGame()

	p1, p2 := igame.GetPlayerIds()

	logging.Info(p1)
	logging.Info(p2)
	igame.MakeMove(p1, 0, 0, 1, 1)
	igame.MakeMove(p2, 7, 7, 6, 6)
}
