package game

import (
	"math"

	"github.com/yelaco/robinhood-chess/internal/config"
	"github.com/yelaco/robinhood-chess/internal/logging"
)

/*
 * Bishop
 */
type bishop struct {
	white         bool
	attachedPiece piece
}

func (b bishop) canMove(board *board, start *spot, end *spot) bool {
	if start == end {
		logging.Info("Same location")
		return false
	} // same location (pointer comparison)

	if end.piece != nil && start.piece.isWhite() == end.piece.isWhite() {
		logging.Info("Same side")
		return false
	} // same side

	if math.Abs(float64(start.x-end.x)) != math.Abs(float64(start.y-end.y)) {
		logging.Info("Invalid move")
		return false
	} // invalid move

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
			logging.Error("Bishop canMove(): Index out of range")
			return false
		}
	}

	return true
}

func (b *bishop) isWhite() bool {
	return b.white
}

func (b *bishop) toUnicode() string {
	if b.white {
		return "♝"
	} else {
		return "♗"
	}
}

func (b *bishop) attach(other piece) {
	if b.attachedPiece != nil {
		logging.Error("The piece is still attaching to other *piece")
	} else {
		b.attachedPiece = other
	}
}

func (b *bishop) detach() piece {
	detachedPiece := b.attachedPiece
	b.attachedPiece = nil
	return detachedPiece
}
