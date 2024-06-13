package game

import (
	"math"

	"github.com/minhquang053/robinhood-chess/internal/logging"
)

/*
 * Knight
 */
type knight struct {
	white         bool
	attachedPiece piece
}

func (k knight) canMove(board *board, start *spot, end *spot) bool {
	if start == end {
		logging.Info("Same location")
		return false
	} // same location (pointer comparison)

	if end.piece != nil && start.piece.isWhite() == end.piece.isWhite() {
		logging.Info("Same side")
		return false
	} // same side

	if math.Abs(float64(start.x-end.x))*math.Abs(float64(start.y-end.y)) != 2.0 {
		logging.Info("Invalid move")
		return false
	} // invalid move

	return true
}

func (k *knight) isWhite() bool {
	return k.white
}

func (k *knight) toUnicode() string {
	if k.white {
		return "♞"
	} else {
		return "♘"
	}
}

func (k *knight) attach(other piece) {
	if k.attachedPiece != nil {
		logging.Error("The piece is still attaching to other *piece")
	} else {
		k.attachedPiece = other
	}
}

func (k *knight) detach() piece {
	detachedPiece := k.attachedPiece
	k.attachedPiece = nil
	return detachedPiece
}
