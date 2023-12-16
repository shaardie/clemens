package pieces

type pieceType uint8

const (
	PawnType pieceType = iota
	BishopType
	KnightType
	RookType
	QueenType
	KingType
)
