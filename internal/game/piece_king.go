package game

import (
	"math"

	"github.com/yelaco/robinhood-chess/pkg/config"
	"github.com/yelaco/robinhood-chess/pkg/logging"
)

/*
 * King
 */
type king struct {
	white     bool
	inCheck   bool
	initMoved bool
}

func (k king) canMove(board *board, start *spot, end *spot) bool {
	if start == end {
		return false
	}

	if end.piece != nil &&
		start.piece.isWhite() == end.piece.isWhite() &&
		!isCastlingMove(start.x, start.y, end.x, end.y) {
		return false
	}

	// check valid move
	if math.Abs(float64(start.x-end.x)) <= 1 && math.Abs(float64(start.y-end.y)) <= 1 {
		if _, ok := end.piece.(*king); ok {
			// if the king can kill the opponent's king, it's already a win, no need to check if being attacked after
			return true
		}

		// simulate that king has taken the other piece
		tempTakenPiece := end.piece
		end.piece = start.piece
		start.piece = nil

		// Restore the board state
		defer func() {
			start.piece = end.piece
			end.piece = tempTakenPiece
		}()

		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				box := board.boxes[i][j]

				// no piece
				if box.piece == nil {
					continue
				}

				// same side
				if box.piece.isWhite() == end.piece.isWhite() {
					continue
				}

				// piece that can attack king
				if box.piece.canMove(board, box, end) {
					return false
				}
			}
		}

		return true
	}

	// if not normal valid move, the only option left is castling
	return k.canCastling(board, start, end)
}

func (k king) isWhite() bool {
	return k.white
}

func (k king) toUnicode() string {
	if k.white {
		return "♚"
	} else {
		return "♔"
	}
}

func isCastlingMove(sx, sy, ex, ey int) bool {
	return sx == 4 && ey == sy && (ex == 0 || ex == 7) && (sy == 0 || sy == 7)
}

func (k *king) canCastling(board *board, start *spot, end *spot) bool {
	if k.initMoved {
		return false
	}

	if r, ok := end.piece.(rook); !ok || r.initMoved {
		return false
	}

	if end.y != start.y {
		return false
	}

	if k.inCheck {
		return false
	}

	if !isCastlingMove(start.x, start.y, end.x, end.y) {
		return false
	}

	inBetweens := make([]*spot, 0, 6)
	inBetweens = append(inBetweens, start)

	i := start.x
	j := start.y
	bex := end.x
	if i < bex {
		bex--
	} else if i > bex {
		bex++
	}

	for i != bex {
		if i < bex {
			i++
		} else {
			i--
		}

		if i > -1 && i < config.BoardLen {
			if board.boxes[i][j].piece != nil {
				return false
			} else {
				inBetweens = append(inBetweens, (board.boxes[i][j]))
			}
		} else {
			logging.Error("King canCastling(): Index out of range")
			return false
		}
	}

	for _, theSpot := range inBetweens {
		for x := 0; x < config.BoardLen; x++ {
			for y := 0; y < config.BoardLen; y++ {
				box := board.boxes[x][y]

				// no piece
				if box.piece == nil {
					continue
				}

				// same side
				if box.piece.isWhite() == start.piece.isWhite() {
					continue
				}

				if box.piece.canMove(board, box, theSpot) {
					return false
				}
			}
		}
	}

	return true
}
