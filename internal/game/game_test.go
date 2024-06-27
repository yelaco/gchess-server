package game

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yelaco/go-chess-server/pkg/utils"
)

func generatePlayerIds() [2]string {
	return [2]string{utils.GenerateUUID(), utils.GenerateUUID()}
}

func setGame(testPiece string) (*Game, string, string) {
	igame := InitGame(generatePlayerIds())
	p1, p2 := igame.GetPlayerIds()
	fmt.Printf("Player 1 ID: %s\nPlayer 2 ID: %s\n\n", p1, p2)
	var moves []string
	switch testPiece {
	case "pawn":
		moves = []string{"e2-e4", "e7-e5", "d2-d4"}
	case "bishop":
		moves = []string{"d2-d4", "e7-e5"}
	case "knight":
		moves = []string{"d2-d4", "e7-e5"}
	case "rook":
		moves = []string{"a2-a4", "b7-b5", "a4-b5", "b8-c6"}
	case "queen":
		moves = []string{"e2-e4", "e7-e5"}
	case "king":
		moves = []string{"e2-e4", "e7-e5", "f1-c4", "f7-f6", "g1-f3", "f8-c5"}
	case "stalemate":
		moves = []string{"e2-e3", "a7-a5", "d1-h5", "a8-a6", "h5-a5", "h7-h5", "h2-h4", "a6-h6",
			"a5-c7", "f7-f6", "c7-d7", "e8-f7", "d7-b7", "d8-d3", "b7-b8", "d3-h7",
			"b8-c8", "f7-g6", "c8-e6"}
	case "checkmate":
		moves = []string{"e2-e4", "e7-e5", "f1-c4", "b8-c6", "d1-h5", "g8-f6", "h5-f7"}
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
	igame, _, p2 := setGame("pawn")

	err := igame.MakeMove(p2, "e5", "d4")
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
		t.Error(err)
	}
	igame.PrintBoard()
}

func TestKing(t *testing.T) {
	igame, p1, _ := setGame("king")
	err := igame.MakeMove(p1, "e1", "h1")
	if err != nil {
		t.Error(err)
	}
	igame.PrintBoard()
}

func TestStalemate(t *testing.T) {
	igame, _, _ := setGame("stalemate")
	if igame.GetStatus() != "STALEMATE" {
		t.Errorf("Test stalemate: got %s, want %s", igame.GetStatus(), "STALEMATE")
	}
	igame.PrintBoard()
}

func TestCheckmate(t *testing.T) {
	igame, _, _ := setGame("checkmate")
	if igame.GetStatus() != "WHITE_CHECKMATE" {
		t.Errorf("Test stalemate: got %s, want %s", igame.GetStatus(), "STALEMATE")
	}
	igame.PrintBoard()
}
