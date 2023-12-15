package bitboard

type BlackPawns Bitboard
type WhitePawns Bitboard

func (p WhitePawns) SinglePushTargets(emptySquares Bitboard) Bitboard {
	return NorthOne(Bitboard(p)) & emptySquares
}

func (p WhitePawns) DoublePushTargets(emptySquares Bitboard) Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(emptySquares)
	// White Double Push only possible on empty fields on rank 4
	return SouthOne(singlePushTargets) & emptySquares & fullRank4
}

func (p WhitePawns) EastAttacks() Bitboard   { return NorthEastOne(Bitboard(p)) }
func (p WhitePawns) WestAttacks() Bitboard   { return NorthWestOne(Bitboard(p)) }
func (p WhitePawns) AnyAttacks() Bitboard    { return p.EastAttacks() | p.WestAttacks() }
func (p WhitePawns) DoubleAttacks() Bitboard { return p.EastAttacks() & p.WestAttacks() }
func (p WhitePawns) SingleAttacks() Bitboard { return p.EastAttacks() ^ p.WestAttacks() }

func (p BlackPawns) SinglePushTargets(emptySquares Bitboard) Bitboard {
	return SouthOne(Bitboard(p)) & emptySquares
}

func (p BlackPawns) DoublePushTargets(emptySquares Bitboard) Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(emptySquares)
	// Black Double Push only possible on empty fields on rank 5
	return SouthOne(singlePushTargets) & emptySquares & fullRank5
}

func (p BlackPawns) EastAttacks() Bitboard   { return SouthEastOne(Bitboard(p)) }
func (p BlackPawns) WestAttacks() Bitboard   { return SouthWestOne(Bitboard(p)) }
func (p BlackPawns) AnyAttacks() Bitboard    { return p.EastAttacks() | p.WestAttacks() }
func (p BlackPawns) DoubleAttacks() Bitboard { return p.EastAttacks() & p.WestAttacks() }
func (p BlackPawns) SingleAttacks() Bitboard { return p.EastAttacks() ^ p.WestAttacks() }

func SafePawnSquares(wp WhitePawns, bp BlackPawns) Bitboard {
	return wp.DoubleAttacks() | ^bp.SingleAttacks() | (wp.SingleAttacks() & ^bp.DoubleAttacks())
}
