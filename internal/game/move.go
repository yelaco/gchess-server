package game

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

func (g *Game) GetLastMove() *move {
	if len(g.moves) == 0 {
		return nil
	}
	lastMove := g.moves[len(g.moves)-1]
	return lastMove
}
