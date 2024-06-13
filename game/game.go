package game

import (
	"errors"
	"fmt"
)

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
		board:       initBoard(),
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

func (g *Game) updateBoard(move *move) {
	switch p := move.pieceMoved.(type) {
	case *pawn:
		p.initMoved = true
	case *king:
		p.initMoved = true
	default:
	}

	if move.isCastling {
	} else if move.isEnpassant {
	} else if move.isPromoting {
	} else {
		move.end.piece = move.pieceMoved
		move.start.piece = nil
	}
}

func (g *Game) MakeMove(playerId, startPos, endPos string) error {
	// check correct turn for move made by player with given id
	if !g.correctTurn(playerId) {
		return fmt.Errorf("wrong turn for player id: %s", playerId)
	}

	// map chess position to board coordinate
	startX, startY := mapChessPosToCoord(startPos)
	endX, endY := mapChessPosToCoord(endPos)

	// setup move
	startBox := g.board.boxes[startX][startY]
	endBox := g.board.boxes[endX][endY]
	move := &move{
		playerId: playerId,
		startPos: startPos,
		endPos:   endPos,
		start:    startBox,
		end:      endBox,
	}

	if err := g.checkMove(move); err != nil {
		return err
	}

	g.updateBoard(move)
	g.checkAndNextTurn()

	// add move to played moves history in the game
	g.moves = append(g.moves, move)

	return nil
}

func (g *Game) checkMove(move *move) error {
	srcPiece := move.start.piece
	if srcPiece == nil {
		return fmt.Errorf("null piece at %d%d", move.start.x, move.start.y)
	}
	dstPiece := move.end.piece

	// check correct turn
	if srcPiece.isWhite() != g.isWhiteTurn {
		return errors.New("wrong turn")
	}

	// check valid move
	if !srcPiece.canMove(g.board, move.start, move.end) {
		return fmt.Errorf("invalid move: %s-%s", move.startPos, move.endPos)
	}

	move.pieceMoved = srcPiece

	if dstPiece != nil {
		move.pieceTaken = dstPiece
	}

	return nil
}

func (g *Game) PrintBoard() {
	fmt.Println("  +-----------------+")

	for i := 7; i >= 0; i-- {
		fmt.Print("  | ")
		for j := 0; j < 8; j++ {
			box := g.board.boxes[j][i]
			if box.piece != nil {
				fmt.Print(box.piece.toUnicode() + " ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println("|")
	}

	fmt.Println("  +-----------------+")
	fmt.Println()
}

func (g *Game) correctTurn(playerId string) bool {
	if g.isWhiteTurn {
		return playerId == g.players[0].playerId
	} else {
		return playerId == g.players[1].playerId
	}
}
