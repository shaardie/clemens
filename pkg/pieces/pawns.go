package pieces

import "github.com/shaardie/clemens/pkg/bitboard"

type BlackPawns bitboard.Bitboard
type WhitePawns bitboard.Bitboard

func (p WhitePawns) SinglePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	return bitboard.NorthOne(bitboard.Bitboard(p)) & emptySquares
}

func (p WhitePawns) DoublePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(emptySquares)
	// White Double Push only possible on empty fields on rank 4
	return bitboard.SouthOne(singlePushTargets) & emptySquares & bitboard.RankMask4
}

func (p WhitePawns) EastAttacks() bitboard.Bitboard {
	return bitboard.NorthEastOne(bitboard.Bitboard(p))
}
func (p WhitePawns) WestAttacks() bitboard.Bitboard {
	return bitboard.NorthWestOne(bitboard.Bitboard(p))
}
func (p WhitePawns) AnyAttacks() bitboard.Bitboard    { return p.EastAttacks() | p.WestAttacks() }
func (p WhitePawns) DoubleAttacks() bitboard.Bitboard { return p.EastAttacks() & p.WestAttacks() }
func (p WhitePawns) SingleAttacks() bitboard.Bitboard { return p.EastAttacks() ^ p.WestAttacks() }

func (p BlackPawns) SinglePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	return bitboard.SouthOne(bitboard.Bitboard(p)) & emptySquares
}

func (p BlackPawns) DoublePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(emptySquares)
	// Black Double Push only possible on empty fields on rank 5
	return bitboard.SouthOne(singlePushTargets) & emptySquares & bitboard.RankMask5
}

func (p BlackPawns) EastAttacks() bitboard.Bitboard {
	return bitboard.SouthEastOne(bitboard.Bitboard(p))
}
func (p BlackPawns) WestAttacks() bitboard.Bitboard {
	return bitboard.SouthWestOne(bitboard.Bitboard(p))
}
func (p BlackPawns) AnyAttacks() bitboard.Bitboard    { return p.EastAttacks() | p.WestAttacks() }
func (p BlackPawns) DoubleAttacks() bitboard.Bitboard { return p.EastAttacks() & p.WestAttacks() }
func (p BlackPawns) SingleAttacks() bitboard.Bitboard { return p.EastAttacks() ^ p.WestAttacks() }

func SafePawnSquares(wp WhitePawns, bp BlackPawns) bitboard.Bitboard {
	return wp.DoubleAttacks() | ^bp.SingleAttacks() | (wp.SingleAttacks() & ^bp.DoubleAttacks())
}
