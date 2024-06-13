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
	err := igame.MakeMove(p1, "c1", "d5")
	if err != nil {
		logging.Info(err.Error())
	}
	logging.Info(igame.GetLastMove())
}
