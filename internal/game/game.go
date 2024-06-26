package game

import (
	"errors"
	"fmt"

	"github.com/yelaco/go-chess-server/pkg/logging"
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
	playerIds   [2]string
	isWhiteTurn bool
	board       *board
	moves       []*move // moves playied through out the game
	status      GameStatus
	kingSpots   [2]*spot
}

func InitGame(playerIds [2]string) *Game {
	g := &Game{
		playerIds:   playerIds,
		board:       initBoard(),
		isWhiteTurn: true,
		status:      active,
	}
	g.kingSpots[0] = g.board.boxes[4][0]
	g.kingSpots[1] = g.board.boxes[4][7]
	return g
}

func (g *Game) GetBoard() [8][8]string {
	boxes := [8][8]string{}
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			box := g.board.boxes[j][i]
			if box.piece != nil {
				boxes[j][i] = box.piece.toUnicode()
			} else {
				boxes[j][i] = ""
			}
		}
	}

	return boxes
}

func (g *Game) GetStatus() string {
	return string(g.status)
}

func (g *Game) GetCurrentTurn() bool {
	return g.isWhiteTurn
}

func (g *Game) GetPlayerSide(playerID string) (bool, error) {
	if g.playerIds[0] == playerID {
		return true, nil
	} else if g.playerIds[1] == playerID {
		return false, nil
	} else {
		return false, errors.New("invalid player id")
	}
}

func (g *Game) GetPlayerIds() (string, string) {
	return g.playerIds[0], g.playerIds[1]
}

func (g *Game) IsOver() bool {
	return g.status != active
}

func (g *Game) checkAndNextTurn(move *move) {
	// go to next turn
	g.isWhiteTurn = !g.isWhiteTurn

	if g.isStalemate() {
		g.status = stalemate
	} else if g.kingInCheck() {
		if g.kingInCheckmate() {
			if g.isWhiteTurn {
				g.status = blackCheckmate
			} else {
				g.status = whiteCheckmate
			}
		}

		move.isChecking = true
	}
}

func (g *Game) updateKingSpots() {
	wk, isWk := g.kingSpots[0].piece.(*king)
	bk, isBk := g.kingSpots[1].piece.(*king)
	if isWk && isBk && wk.isWhite() && !bk.isWhite() {
		return
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if k, ok := g.board.boxes[x][y].piece.(*king); ok {
				if k.isWhite() {
					g.kingSpots[0] = g.board.boxes[x][y]
				} else {
					g.kingSpots[1] = g.board.boxes[x][y]
				}
			}
		}
	}
}

func (g *Game) kingInCheck() bool {
	g.updateKingSpots()
	var kingSpot *spot
	if g.isWhiteTurn {
		kingSpot = g.kingSpots[0]
	} else {
		kingSpot = g.kingSpots[1]
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			box := g.board.boxes[i][j]
			if box.piece == nil {
				continue
			}
			if box.piece.isWhite() == g.isWhiteTurn {
				continue
			}
			if box.piece.canMove(g.board, box, kingSpot) {
				if k, ok := kingSpot.piece.(*king); ok {
					k.inCheck = true
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) kingInCheckmate() bool {
	threatPieceSpots := []*spot{}
	var kingSpot *spot
	if g.isWhiteTurn {
		kingSpot = g.kingSpots[0]
	} else {
		kingSpot = g.kingSpots[1]
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			box := g.board.boxes[i][j]
			if box.piece == nil {
				continue
			}
			if box.piece.isWhite() == g.isWhiteTurn {
				continue
			}
			if box.piece.canMove(g.board, box, kingSpot) {
				threatPieceSpots = append(threatPieceSpots, box)
			}
		}
	}

	king, ok := kingSpot.piece.(*king)
	if !ok {
		logging.Error("kingInCheckmate(): king not in spot")
	}

	// if there are multiple threats, the King must run (can't block attack)
	// we check if the King can run or not with any number of threats
	for i := kingSpot.x - 1; i <= kingSpot.x+1; i++ {
		for j := kingSpot.y - 1; j <= kingSpot.y+1; j++ {
			if i < 0 || i > 7 || j < 0 || j > 7 {
				continue
			}
			if king.canMove(g.board, kingSpot, g.board.boxes[i][j]) {
				return false
			}
		}
	}

	// the king can't run and there are multiple checking pieces
	// ==> can't block
	if len(threatPieceSpots) > 1 {
		return true
	}

	threatSpot := threatPieceSpots[0]

	// if there is a single threat and the King can't run, the threatening piece must be stopped
	// first by eliminate the threatening piece
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			box := g.board.boxes[i][j]
			if box.piece == nil {
				continue
			}
			if box.piece.isWhite() != g.isWhiteTurn {
				continue
			}
			if box.piece.canMove(g.board, box, threatSpot) {
				return false
			}
		}
	}

	// second by block the attack
	// with knight or pawn checking, there aren't any spots in between to block the attack
	if _, ok := threatSpot.piece.(*knight); ok {
		return true
	}
	if _, ok := threatSpot.piece.(*pawn); ok {
		return true
	}

	// get pieces in between the threat piece and the king
	inBetweens := make([]*spot, 0, 6)
	i := threatSpot.x
	j := threatSpot.y
	bex := kingSpot.x
	bey := kingSpot.y
	if i < bex {
		bex--
	} else if i > bex {
		bex++
	}
	if j < bey {
		bey--
	} else if j > bey {
		bey++
	}

	for i != bex || j != bey {
		if i < bex {
			i++
		} else if i > bex {
			i--
		}
		if j < bey {
			j++
		} else if j > bey {
			j--
		}

		inBetweens = append(inBetweens, g.board.boxes[i][j])
	}

	// check if the check can be blocked
	for _, theSpot := range inBetweens {
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				box := g.board.boxes[x][y]
				if box.piece == nil {
					continue
				}
				if box.piece.isWhite() != king.isWhite() {
					continue
				}
				if box.piece.canMove(g.board, box, theSpot) {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) isStalemate() bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			box := g.board.boxes[i][j]
			if box.piece == nil {
				continue
			}
			if box.piece.isWhite() != g.isWhiteTurn {
				continue
			}

			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					if !box.piece.canMove(g.board, box, g.board.boxes[x][y]) {
						continue
					}

					tmpPiece := g.board.boxes[x][y].piece
					g.board.boxes[x][y].piece = box.piece
					box.piece = nil

					if !g.kingInCheck() {
						box.piece = g.board.boxes[x][y].piece
						g.board.boxes[x][y].piece = tmpPiece
						return false
					}

					box.piece = g.board.boxes[x][y].piece
					g.board.boxes[x][y].piece = tmpPiece
				}
			}
		}
	}

	return true
}

func (g *Game) updateBoard(move *move) {
	move.start.piece = nil
	switch p := move.pieceMoved.(type) {
	case *pawn:
		move.end.piece = move.pieceMoved
		if !p.initMoved {
			move.isInitMove = true
			p.initMoved = true
		}
		if move.isEnpassant {
			move.end.piece = move.pieceMoved
			g.board.boxes[move.end.x][move.start.y].piece = nil
		} else if move.isPromoting {
			move.end.piece = p.promote("queen")
			move.piecePromoted = move.end.piece
		}
	case *king:
		if !p.initMoved {
			move.isInitMove = true
			p.initMoved = true
		}
		if p.isWhite() {
			if move.isCastling {
				if move.end.x == 0 {
					g.kingSpots[0] = g.board.boxes[2][0]
					g.board.boxes[3][0].piece = move.pieceTaken
				} else if move.end.x == 7 {
					g.kingSpots[0] = g.board.boxes[6][0]
					g.board.boxes[5][0].piece = move.pieceTaken
				}
				g.kingSpots[0].piece = move.pieceMoved
				move.end.piece = nil
			} else {
				g.kingSpots[0] = move.end
			}
		} else {
			if move.isCastling {
				if move.end.x == 0 {
					g.kingSpots[1] = g.board.boxes[2][7]
					g.board.boxes[3][7].piece = move.pieceTaken
				} else if move.end.x == 7 {
					g.kingSpots[1] = g.board.boxes[6][7]
					g.board.boxes[5][7].piece = move.pieceTaken
				}
				g.kingSpots[1].piece = move.pieceMoved
				move.end.piece = nil
			} else {
				move.end.piece = move.pieceMoved
				g.kingSpots[0] = move.end
			}
		}
	case *rook:
		if !p.initMoved {
			move.isInitMove = true
			p.initMoved = true
		}
		move.end.piece = move.pieceMoved
	default:
		move.end.piece = move.pieceMoved
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
	g.checkAndNextTurn(move)

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
		return errors.New("can't play your opponent's piece")
	}

	// check valid move
	switch p := srcPiece.(type) {
	case *pawn:
		if !p.canMove(g.board, move.start, move.end) {
			if p.canEnpassant(move.start, move.end, g.GetLastMove()) {
				move.isEnpassant = true
			} else {
				return fmt.Errorf("invalid pawn move: %s-%s", move.startPos, move.endPos)
			}
		}
		if move.end.y == 7 || move.end.y == 0 {
			move.isPromoting = true
		}
	case *king:
		move.isCastling = isCastlingMove(move.start.x, move.start.y, move.end.x, move.end.y)
		if !p.canMove(g.board, move.start, move.end) {
			return fmt.Errorf("invalid king move: %s-%s", move.startPos, move.endPos)
		}
	default:
		if !p.canMove(g.board, move.start, move.end) {
			return fmt.Errorf("invalid move: %s-%s", move.startPos, move.endPos)
		}
	}

	if !move.isCastling {
		move.end.piece = srcPiece
		move.start.piece = nil
		if g.kingInCheck() {
			move.start.piece = srcPiece
			move.end.piece = dstPiece
			return fmt.Errorf("invalid move: %s-%s, king in checked", move.startPos, move.endPos)
		}
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
		return playerId == g.playerIds[0]
	} else {
		return playerId == g.playerIds[1]
	}
}
