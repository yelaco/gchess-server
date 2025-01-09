package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.NewString()
}

func BoardToFen(board [8][8]string) string {
	var fen strings.Builder

	for j := 7; j >= 0; j-- {
		emptyCount := 0
		for i := range 8 {
			if board[i][j] == "." || board[i][j] == "" { // Treat "." or "" as empty square
				emptyCount++
			} else {
				if emptyCount > 0 {
					fen.WriteString(strconv.Itoa(emptyCount))
					emptyCount = 0
				}
				switch board[i][j] {
				case "♖":
					fen.WriteString("r")
				case "♘":
					fen.WriteString("n")
				case "♗":
					fen.WriteString("b")
				case "♕":
					fen.WriteString("q")
				case "♔":
					fen.WriteString("k")
				case "♙":
					fen.WriteString("p")
				case "♜":
					fen.WriteString("R")
				case "♞":
					fen.WriteString("N")
				case "♝":
					fen.WriteString("B")
				case "♛":
					fen.WriteString("Q")
				case "♚":
					fen.WriteString("K")
				case "♟":
					fen.WriteString("P")
				default:
					panic("Invalid piece symbol: " + board[i][j])
				}
			}
		}
		if emptyCount > 0 {
			fen.WriteString(strconv.Itoa(emptyCount))
		}
		fen.WriteString("/") // Add row separator
	}

	fenStr := fen.String()
	return fenStr[:len(fenStr)-1] // Remove trailing "/"
}
