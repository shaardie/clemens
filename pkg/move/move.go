package move

import "github.com/shaardie/clemens/pkg/types"

type MoveType int

const (
	NORMAL MoveType = iota
	PROMOTION
	EN_PASSANT
	CASTLING
)

// Move represents a move from one position to another
// 0-5 is the source square
// 6-11 is the destination square
// 12-13 is the Move Type
// 14-15 is the Promotion Piece Type
type Move uint64

func (m *Move) GetSourceSquare() int {
	return int(*m & 0b111111)
}

func (m *Move) SetSourceSquare(square int) {
	*m |= Move(square)
}

func (m *Move) GetDestinationSquare() int {
	return int(*m >> 6 & 0b111111)
}

func (m *Move) SetDestinationSquare(square int) {
	*m |= Move(square << 6)
}

func (m *Move) GetMoveType() MoveType {
	return MoveType(*m >> 12 & 0b11)
}

func (m *Move) SetMoveType(mt MoveType) {
	*m |= Move(mt << 12)
}

func (m *Move) GetPromitionPieceType() types.PieceType {
	return types.PieceType(*m>>14) + 1
}

func (m *Move) SetPromitionPieceType(pt types.PieceType) {
	*m |= Move((pt - 1) << 14)
}
