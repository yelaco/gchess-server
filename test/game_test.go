package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yelaco/robinhood-chess/internal/game"
)

func setGame(testPiece string) (*game.Game, string, string) {
	igame := game.InitGame()
	p1, p2 := igame.GetPlayerIds()
	fmt.Printf("Player 1 ID: %s\nPlayer 2 ID: %s\n\n", p1, p2)
	var moves []string
	switch testPiece {
	case "pawn":
		moves = []string{"d2-d4", "b8-c6", "d4-d5", "e7-e5", "d5-e6", "a7-a6", "e6-d7", "d8-d7"}
	case "bishop":
		moves = []string{"d2-d4", "e7-e5"}
	case "knight":
		moves = []string{"d2-d4", "e7-e5"}
	case "rook":
		moves = []string{"a2-a4", "b7-b5", "a4-b5", "b8-c6"}
	case "queen":
		moves = []string{"e2-e4", "e7-e5"}
	case "king":
		moves = []string{"e2-e4", "e7-e5", "f1-c4", "f7-f6", "g1-f3", "f8-c5", "e1-h1", "g8-h6", "f1-e1", "d7-d5", "a2-a3", "e8-h8"}
	default:
	}
	for i, move := range moves {
		pos := strings.Split(move, "-")
		if i%2 == 0 {
			err := igame.MakeMove(p1, pos[0], pos[1])
			if err != nil {
				fmt.Printf("%s: %s\n", move, err.Error())
			}
		} else {
			err := igame.MakeMove(p2, pos[0], pos[1])
			if err != nil {
				fmt.Printf("%s: %s\n", move, err.Error())
			}
		}

		// igame.PrintBoard()
	}
	return igame, p1, p2
}

func TestPawn(t *testing.T) {
	igame, p1, _ := setGame("pawn")

	// err := igame.MakeMove(p1, "d4", "d6")
	// err := igame.MakeMove(p1, "d4", "e5")
	err := igame.MakeMove(p1, "d7", "c8")
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
	}
	igame.PrintBoard()
}

func TestRook(t *testing.T) {
	igame, p1, _ := setGame("rook")
	err := igame.MakeMove(p1, "a1", "a7")
	if err != nil {
		fmt.Println(err.Error())
	}
	igame.PrintBoard()
}

func TestQueen(t *testing.T) {
	igame, p1, _ := setGame("queen")
	err := igame.MakeMove(p1, "d1", "g4")
	if err != nil {
		fmt.Println(err.Error())
	}
	igame.PrintBoard()
}

func TestKing(t *testing.T) {
	igame, p1, _ := setGame("king")
	err := igame.MakeMove(p1, "e1", "h1")
	if err != nil {
		fmt.Println(err.Error())
	}
	igame.PrintBoard()
}
