package game

import (
	"math"

	"github.com/minhquang053/robinhood-chess/internal/config"
	"github.com/minhquang053/robinhood-chess/internal/logging"
)

type piece interface {
	canMove(b board, start spot, end spot) bool
	isWhite() bool
	attach(*piece)
	// detach()
}

/*
 * Bishop
 */
type bishop struct {
	white         bool
	attachedPiece *piece
}

func (b bishop) canMove(board board, start spot, end spot) bool {
	if (start.x == end.x && start.y == end.y) || // same location
		(end.piece != nil && start.piece.isWhite() == end.piece.isWhite()) || // same side
		(math.Abs(float64(start.x-end.x)) != math.Abs(float64(start.y-end.y))) { // invalid move
		return false
	}

	i := start.x
	j := start.y
	bex := end.x // bound of end x to check starting from i
	bey := end.y // bound of end y to check starting from j
	if i < bex {
		bex--
	} else {
		bex++
	}
	if j < bey {
		bey--
	} else {
		bey++
	}

	for i != bex && j != bey {
		if i < bex {
			i++
		} else {
			i--
		}
		if j < bey {
			j++
		} else {
			j--
		}

		if i < config.BoardLen && j < config.BoardLen {
			if board.boxes[i][j].piece != nil {
				return false
			}
		} else {
			logging.Panic("Bishop canMove(): Index out of range")
		}
	}

	return true
}

func (b *bishop) isWhite() bool {
	return b.white
}

func (b *bishop) attach(other *piece) {
	if b.attachedPiece != nil {
		logging.Error("The piece is still attaching to other *piece")
	} else {
		b.attachedPiece = other
	}
}

/*
 * Knight
 */
type knight struct {
	white         bool
	attachedPiece piece
}

func (k knight) canMove(board board, start spot, end spot) bool {

	return true
}

func (k *knight) isWhite() bool {
	return k.white
}

func (k *knight) attach(other *piece) {
	// Implement the actual logic for attaching a piece
}

/*
 * Rook
 */
type rook struct {
	white         bool
	attachedPiece piece
}

func (r rook) canMove(board board, start spot, end spot) bool {
	return true
}

func (r *rook) isWhite() bool {
	return r.white
}

func (r *rook) attach(other *piece) {

}

/*
 * Queen
 */
type queen struct {
	white         bool
	attachedPiece piece
}

func (q queen) canMove(b board, start spot, end spot) bool {
	return true
}

func (q *queen) isWhite() bool {
	return q.white
}

func (q *queen) attach(other *piece) {
	// Implement the actual logic for attaching a piece
}

/*
 * Pawn
 */
type pawn struct {
	white         bool
	attachedPiece piece
}

func (p pawn) canMove(b board, start spot, end spot) bool {
	return true
}

func (p *pawn) isWhite() bool {
	return p.white
}

func (p *pawn) attach(other *piece) {
	// Implement the actual logic for attaching a piece
}

/*
 * King
 */
type king struct {
	white         bool
	attachedPiece piece
}

func (k king) canMove(b board, start spot, end spot) bool {
	return true
}

func (k *king) isWhite() bool {
	return k.white
}

func (k *king) attach(other *piece) {
	// Implement the actual logic for attaching a piece
}
