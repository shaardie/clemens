package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	isolanis = [game_number]int16{-20, -5}

	passedScalar int16 = 5
)

// evalPawns evaluates the pawn structure
func (e *eval) evalPawns(pos *position.Position) {
	whitePawns := pos.PiecesBitboard[types.WHITE][types.PAWN]
	blackPawns := pos.PiecesBitboard[types.BLACK][types.PAWN]

	isolaniDiff := int16(pawn.NumberOfIsolanis(whitePawns) - pawn.NumberOfIsolanis(blackPawns))
	for phase := range game_number {
		e.phaseScores[phase] += isolanis[phase] * isolaniDiff
	}

	e.baseScore += passed(whitePawns, blackPawns)
}

func passed(whitePawns, blackPawns bitboard.Bitboard) int16 {
	passedWhite := pawn.Passed(types.WHITE, whitePawns, blackPawns)
	passedBlack := pawn.Passed(types.BLACK, whitePawns, blackPawns)
	return rankedPawnEval(passedScalar, passedWhite, passedBlack)
}

func rankedPawnEval(scalar int16, selectedWhitePawns, selectedBlackPawns bitboard.Bitboard) int16 {
	var r int16
	rankMask := bitboard.RankMask2
	for rank := types.RANK_2; rank <= types.RANK_7; rank++ {
		whiteFactor := int(rank)
		blackFactor := int(7 - rank)
		r = r + scalar*int16(whiteFactor*(selectedWhitePawns&rankMask).PopulationCount()-
			blackFactor*(selectedBlackPawns&rankMask).PopulationCount())
		rankMask = rankMask << 8
	}
	return r
}
