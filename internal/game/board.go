package game

import "github.com/yelaco/robinhood-chess/pkg/config"

type spot struct {
	piece piece
	x     int
	y     int
}

type board struct {
	boxes [8][8]*spot
}

func initBoard() *board {
	b := &board{}

	// Set spots for each box in the board
	for i := range b.boxes {
		for j := range b.boxes[i] {
			b.boxes[i][j] = &spot{
				piece: nil,
				x:     i,
				y:     j,
			}
		}
	}

	// Set pieces to their initial positions on the board
	b.boxes[0][0].piece = rook{white: true}
	b.boxes[1][0].piece = knight{white: true}
	b.boxes[2][0].piece = bishop{white: true}
	b.boxes[3][0].piece = queen{white: true}
	b.boxes[4][0].piece = king{white: true}
	b.boxes[5][0].piece = bishop{white: true}
	b.boxes[6][0].piece = knight{white: true}
	b.boxes[7][0].piece = rook{white: true}

	b.boxes[0][7].piece = rook{white: false}
	b.boxes[1][7].piece = knight{white: false}
	b.boxes[2][7].piece = bishop{white: false}
	b.boxes[3][7].piece = queen{white: false}
	b.boxes[4][7].piece = king{white: false}
	b.boxes[5][7].piece = bishop{white: false}
	b.boxes[6][7].piece = knight{white: false}
	b.boxes[7][7].piece = rook{white: false}

	for i := range config.BoardLen {
		b.boxes[i][1].piece = pawn{white: true}
		b.boxes[i][6].piece = pawn{white: false}
	}

	return b
}
