package game

type GameStatus string

const (
	active         GameStatus = "ACTIVE"
	blackCheckmate GameStatus = "BLACK_CHECKMATE"
	whiteCheckmate GameStatus = "WHITE_CHECKMATE"
	stalemate      GameStatus = "STALEMATE"
	blackResign    GameStatus = "BLACK_RESIGN"
	whiteResign    GameStatus = "WHITE_RESIGN"
)

type Game struct {
	players     *[2]player
	isWhiteTurn bool
	board       *board
	moves       []*move // moves played through out the game
	status      GameStatus
}

func InitGame() *Game {
	return &Game{
		players:     initPlayers(),
		board:       &board{},
		isWhiteTurn: true,
		status:      active,
	}
}

func (g *Game) GetPlayerIds() (string, string) {
	return g.players[0].playerId, g.players[1].playerId
}

func (g *Game) IsOver() bool {
	return g.status != active
}

func (g *Game) checkAndNextTurn() {

	// go to next turn
	g.isWhiteTurn = !g.isWhiteTurn
}

func (g *Game) MakeMove(playerId string, startX int, startY int, endX int, endY int) MoveStatus {
	// setup move
	startBox := g.board.boxes[startX][startY]
	endBox := g.board.boxes[endX][endY]
	move := &move{
		playerId: playerId,
		start:    startBox,
		end:      endBox,
	}

	if moveStatus := g.checkMove(move); moveStatus != successMove {
		return moveStatus
	}

	g.checkAndNextTurn()

	// add move to played moves history in the game
	g.moves = append(g.moves, move)
	return successMove
}

func (g *Game) checkMove(move *move) MoveStatus {
	srcPiece := move.start.piece
	if srcPiece == nil {
		return nullPieceMove
	}
	dstPiece := move.end.piece

	// check correct turn
	if srcPiece.isWhite() != g.isWhiteTurn {
		return wrongTurnMove
	}

	// check valid move
	if !srcPiece.canMove(g.board, &move.start, &move.end) {
		return cannotMove
	}

	move.pieceMoved = srcPiece

	if dstPiece != nil {
		move.pieceTaken = dstPiece
	}

	return successMove
}
