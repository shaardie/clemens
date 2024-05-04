package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

type (
	shieldSide int
	shieldRank int
)

const (
	shieldSideKing shieldSide = iota
	shieldSideQueen
	shieldSideNumber
)

const (
	shieldRank2 shieldRank = iota
	shieldRank3
	shieldRankNumber
)

var shieldValue = [shieldRankNumber]int16{40, 20}

var pawnShield = [types.COLOR_NUMBER][shieldSideNumber][shieldRankNumber]bitboard.Bitboard{}

func init() {
	pawnShield[types.WHITE][shieldSideQueen][shieldRank2] = bitboard.BitBySquares(
		types.SQUARE_A2, types.SQUARE_B2, types.SQUARE_C2,
	)
	pawnShield[types.WHITE][shieldSideQueen][shieldRank3] = bitboard.BitBySquares(
		types.SQUARE_A3, types.SQUARE_B3, types.SQUARE_C3,
	)
	pawnShield[types.WHITE][shieldSideKing][shieldRank2] = bitboard.BitBySquares(
		types.SQUARE_F2, types.SQUARE_G2, types.SQUARE_H2,
	)
	pawnShield[types.WHITE][shieldSideKing][shieldRank3] = bitboard.BitBySquares(
		types.SQUARE_F3, types.SQUARE_G3, types.SQUARE_H3,
	)
	pawnShield[types.BLACK][shieldSideQueen][shieldRank2] = bitboard.BitBySquares(
		types.SQUARE_A7, types.SQUARE_B7, types.SQUARE_C7,
	)
	pawnShield[types.BLACK][shieldSideQueen][shieldRank3] = bitboard.BitBySquares(
		types.SQUARE_A6, types.SQUARE_B6, types.SQUARE_C6,
	)
	pawnShield[types.BLACK][shieldSideKing][shieldRank2] = bitboard.BitBySquares(
		types.SQUARE_F7, types.SQUARE_G7, types.SQUARE_H7,
	)
	pawnShield[types.BLACK][shieldSideKing][shieldRank3] = bitboard.BitBySquares(
		types.SQUARE_F6, types.SQUARE_G6, types.SQUARE_H6,
	)

}

func (e *eval) evalKingShield(pos *position.Position) {
	e.phaseScores[midgame] += kingShield(pos)
}

func kingShield(pos *position.Position) int16 {
	var r int16
	// White
	kingSquare := bitboard.LeastSignificantOneBit(pos.PiecesBitboard[types.WHITE][types.KING])
	if types.RankOfSquare(kingSquare) == types.RANK_1 {
		r += kingShieldValue(types.FileOfSquare(kingSquare), types.WHITE, pos.PiecesBitboard[types.WHITE][types.PAWN])
	}

	// Black
	kingSquare = bitboard.LeastSignificantOneBit(pos.PiecesBitboard[types.BLACK][types.KING])
	if types.RankOfSquare(kingSquare) == types.RANK_8 {
		r -= kingShieldValue(types.FileOfSquare(kingSquare), types.BLACK, pos.PiecesBitboard[types.BLACK][types.PAWN])
	}
	return r
}

func kingShieldValue(kingFile uint8, c types.Color, p bitboard.Bitboard) int16 {
	var s shieldSide
	if kingFile > types.FILE_E {
		s = shieldSideKing
	} else if kingFile < types.FILE_D {
		s = shieldSideQueen
	} else {
		// No Bonus
		return 0
	}

	var r int16
	for k := range shieldRankNumber {
		r += shieldValue[k] * int16((pawnShield[c][s][k] & p).PopulationCount())
	}
	return r
}
