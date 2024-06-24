package game

import (
	"math"

	"github.com/yelaco/robinhood-chess/pkg/config"
	"github.com/yelaco/robinhood-chess/pkg/logging"
)

/*
 * Pawn
 */
type pawn struct {
	white         bool
	attachedPiece piece
	initMoved     bool // 2 step init move
}

func (p pawn) canMove(board *board, start *spot, end *spot) bool {
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
		if math.Abs(float64(end.y-start.y)) == 2.0 {
			direction := (end.y - start.y) / 2

			if start.y+direction < config.BoardLen {
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
		if math.Abs(float64(end.y-start.y)) == 2.0 {
			return false
		}
		return end.piece != nil
	}
}

func (p pawn) canEnpassant(start *spot, end *spot, lastMove *move) bool {
	if lastMove == nil {
		return false
	}
	_, ok := lastMove.pieceMoved.(pawn)
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
		logging.Fatal("Pawn promoted to UNDEFINED")
		return nil
	}
}

func (p *pawn) attach(other piece) {
	if p.attachedPiece != nil {
		logging.Error("The piece is still attaching to other *piece")
	} else {
		p.attachedPiece = other
	}
}

func (p *pawn) detach() piece {
	detachedPiece := p.attachedPiece
	p.attachedPiece = nil
	return detachedPiece
}
