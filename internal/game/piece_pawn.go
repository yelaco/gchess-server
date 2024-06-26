package game

import (
	"math"

	"github.com/yelaco/go-chess-server/pkg/logging"
)

/*
 * Pawn
 */
type pawn struct {
	white     bool
	initMoved bool // 2 step init move
}

func (p *pawn) canMove(board *board, start *spot, end *spot) bool {
	if start == end {
		return false
	} // same location (pointer comparison)

	if end.piece != nil && start.piece.isWhite() == end.piece.isWhite() {
		return false
	} // same side

	if (start.piece.isWhite() && start.y > end.y) ||
		(!start.piece.isWhite() && start.y < end.y) {
		return false
	}

	if math.Abs(float64(end.x-start.x)) > 1.0 ||
		math.Abs(float64(end.y-start.y)) > 2.0 ||
		math.Abs(float64(end.y-start.y)) == 0.0 {
		return false
	} // invalid move

	if end.x == start.x {
		if end.piece != nil {
			return false
		}
		if direction := end.y - start.y; direction == 2 || direction == -2 {
			direction /= 2

			if start.y+direction < 8 {
				if board.boxes[start.x][start.y+direction].piece != nil {
					return false
				}
			} else {
				logging.Error("Pawn canMove(): Index out of range")
			}
			return !p.initMoved
		}
		return true
	} else {
		if direction := end.y - start.y; direction == 2 || direction == -2 {
			return false
		}
		return end.piece != nil
	}
}

func (p pawn) canEnpassant(start *spot, end *spot, lastMove *move) bool {
	if lastMove == nil {
		return false
	}
	_, ok := lastMove.pieceMoved.(*pawn)
	return ok && lastMove.isInitMove && lastMove.end.y == start.y && lastMove.end.x == end.x
}

func (p pawn) isWhite() bool {
	return p.white
}

func (p pawn) toUnicode() string {
	if p.white {
		return "♟"
	} else {
		return "♙"
	}
}

func (p pawn) promote(pieceName string) piece {
	switch pieceName {
	case "bishop":
		return &bishop{white: p.white}
	case "knight":
		return &knight{white: p.white}
	case "rook":
		return &rook{white: p.white}
	case "queen":
		return &queen{white: p.white}
	default:
		logging.Error("Pawn promoted to UNDEFINED")
		return nil
	}
}
