package bitboard

type Pawns interface {
	SinglePushTargets(Pawns, emptySquares Bitboard) Bitboard
	DoublePushTargets(Pawns, emptySquares Bitboard) Bitboard
	AbleToSinglePush(Pawns, emptySquares Bitboard) Bitboard
	AbleToDoublePush(Pawns, emptySquares Bitboard) Bitboard
}

type BlackPawns Bitboard
type WhitePawns Bitboard

func (p WhitePawns) SinglePushTargets(Pawns, emptySquares Bitboard) Bitboard {
	return NorthOne(Pawns) & emptySquares
}

func (p WhitePawns) DoublePushTargets(Pawns, emptySquares Bitboard) Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(Pawns, emptySquares)
	// White Double Push only possible on empty fields on rank 4
	return SouthOne(singlePushTargets) & emptySquares & rank4
}

func (p BlackPawns) SinglePushTargets(Pawns, emptySquares Bitboard) Bitboard {
	return SouthOne(Pawns) & emptySquares
}

func (p BlackPawns) DoublePushTargets(Pawns, emptySquares Bitboard) Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(Pawns, emptySquares)
	// Black Double Push only possible on empty fields on rank 5
	return SouthOne(singlePushTargets) & emptySquares & rank5
}
