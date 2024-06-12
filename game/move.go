package game

type MoveStatus string

const (
	successMove   MoveStatus = "SUCCESS"
	wrongTurnMove MoveStatus = "WRONG_TURN"
	cannotMove    MoveStatus = "CANT_MOVE"
	nullPieceMove MoveStatus = "NULL_PIECE"
)

type move struct {
	playerId      string
	start         spot
	end           spot
	pieceMoved    piece
	pieceTaken    piece
	piecePromoted piece
	isCastling    bool
	isChecking    bool
	isEnpassant   bool
	isPromoting   bool
}
