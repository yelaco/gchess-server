package game

import (
	"errors"
	"strings"
)

type MoveStatus string

type move struct {
	playerId      string
	startPos      string
	endPos        string
	start         *spot
	end           *spot
	pieceMoved    piece
	pieceTaken    piece
	piecePromoted piece
	isCastling    bool
	isChecking    bool
	isEnpassant   bool
	isPromoting   bool
	isInitMove    bool
}

func mapChessPosToCoord(pos string) (x int, y int) {
	x = int(pos[0] - 'a')
	y = int(pos[1] - '1')
	return
}

func ParseMove(move string) ([]string, error) {
	move = strings.TrimSpace(move)
	if len(move) != 5 {
		return []string{}, errors.New("couldn't parse move")
	}

	if move[2] != '-' {
		return []string{}, errors.New("couldn't parse move")
	}

	if move[0] >= 'a' && move[0] <= 'h' &&
		move[1] >= '1' && move[1] <= '8' &&
		move[3] >= 'a' && move[3] <= 'h' &&
		move[4] >= '1' && move[4] <= '8' {
		pos := strings.Split(move, "-")
		return pos, nil
	} else {
		return []string{}, errors.New("couldn't parse move")
	}
}
func IsValidMove(move string) bool {
	move = strings.TrimSpace(move)
	if len(move) != 5 {
		return false
	}

	if move[2] != '-' {
		return false
	}

	return move[0] >= 'a' && move[0] <= 'h' &&
		move[1] >= '1' && move[1] <= '8' &&
		move[3] >= 'a' && move[3] <= 'h' &&
		move[4] >= '1' && move[3] <= '8'
}

func (g *Game) GetLastMove() *move {
	if len(g.moves) == 0 {
		return nil
	}
	lastMove := g.moves[len(g.moves)-1]
	return lastMove
}

func (g *Game) GetAllMoves() []string {
	res := make([]string, 0, len(g.moves))
	for _, move := range g.moves {
		res = append(res, move.startPos+"-"+move.endPos)
	}
	return res
}
