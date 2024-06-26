package game

type piece interface {
	canMove(b *board, start *spot, end *spot) bool
	isWhite() bool
	toUnicode() string
}
