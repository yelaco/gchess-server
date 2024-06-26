package game

import (
	"math"
)

/*
 * Knight
 */
type knight struct {
	white bool
}

func (k *knight) canMove(board *board, start *spot, end *spot) bool {
	if start == end {
		return false
	} // same location (pointer comparison)

	if end.piece != nil && start.piece.isWhite() == end.piece.isWhite() {
		return false
	} // same side

	if math.Abs(float64(start.x-end.x))*math.Abs(float64(start.y-end.y)) != 2.0 {
		return false
	} // invalid move

	return true
}

func (k knight) isWhite() bool {
	return k.white
}

func (k knight) toUnicode() string {
	if k.white {
		return "♞"
	} else {
		return "♘"
	}
}
