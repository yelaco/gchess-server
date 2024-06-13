package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/minhquang053/robinhood-chess/game"
)

func setGame(testPiece string) (*game.Game, string, string) {
	igame := game.InitGame()
	p1, p2 := igame.GetPlayerIds()
	fmt.Printf("Player 1 ID: %s\nPlayer2 ID: %s\n\n", p1, p2)
	var moves []string
	switch testPiece {
	case "pawn":
		moves = []string{"d2-d4", "e7-e5"}
	case "bishop":
		moves = []string{"d2-d4", "e7-e5"}
	case "knight":
		moves = []string{"d2-d4", "e7-e5"}
	default:
	}
	for i, move := range moves {
		pos := strings.Split(move, "-")
		if i%2 == 0 {
			err := igame.MakeMove(p1, pos[0], pos[1])
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			err := igame.MakeMove(p2, pos[0], pos[1])
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		// igame.PrintBoard()
	}
	return igame, p1, p2
}

func TestPawn(t *testing.T) {
	igame, p1, _ := setGame("pawn")

	// err := igame.MakeMove(p1, "d4", "d6")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	err := igame.MakeMove(p1, "d4", "e5")
	if err != nil {
		fmt.Println(err.Error())
	}

	igame.PrintBoard()
}

func TestBishop(t *testing.T) {
	igame, p1, _ := setGame("bishop")

	err := igame.MakeMove(p1, "c1", "f4")
	if err != nil {
		fmt.Println(err.Error())
	}
	igame.PrintBoard()
}

func TestKnight(t *testing.T) {
	igame, p1, _ := setGame("knight")
	err := igame.MakeMove(p1, "b1", "d2")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Last move: %s\n", igame.GetLastMove())
	}
	igame.PrintBoard()
}
