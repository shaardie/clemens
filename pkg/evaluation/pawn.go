package evaluation

import (
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	isolanis  = [game_number]int16{-20, -5}
	doubled   = [game_number]int16{-5, -15}
	passPawns = [game_number]int16{20, 80}
	backwards = [game_number]int16{0, 0}
	supported = [game_number]int16{0, 0}
	phalanx   = [game_number]int16{0, 0}
	opposed   = [game_number]int16{0, 0}
)

// evalPawns evaluates the pawn structure
func (e *eval) evalPawns(pos *position.Position) {
	whitePawns := pos.PiecesBitboard[types.WHITE][types.PAWN]
	blackPawns := pos.PiecesBitboard[types.BLACK][types.PAWN]
	opposedWhite := pawn.Opposed(types.WHITE, whitePawns, blackPawns)
	opposedBlack := pawn.Opposed(types.BLACK, whitePawns, blackPawns)
	phalanxWhite := pawn.Phalanx(whitePawns)
	phalanxBlack := pawn.Phalanx(blackPawns)
	supportedWhite := pawn.Supported(types.WHITE, whitePawns)
	supportedBlack := pawn.Supported(types.WHITE, whitePawns)

	isolaniDiff := int16(pawn.NumberOfIsolanis(whitePawns) - pawn.NumberOfIsolanis(blackPawns))
	doublePawnDiff := int16(pawn.NumberOfDoubled(whitePawns) - pawn.NumberOfDoubled(blackPawns))
	passedPawnDiff := int16(pawn.NumberOfPassed(types.WHITE, whitePawns, blackPawns) - pawn.NumberOfPassed(types.BLACK, whitePawns, blackPawns))
	backwardsPawnDiff := int16(pawn.NumberOfBackwards(types.WHITE, whitePawns, blackPawns) - pawn.NumberOfBackwards(types.BLACK, whitePawns, blackPawns))
	opposedDiff := int16(opposedWhite.PopulationCount() - opposedBlack.PopulationCount())
	phalanxDiff := int16(phalanxWhite.PopulationCount() - phalanxBlack.PopulationCount())
	supportedDiff := int16(supportedWhite.PopulationCount() - supportedBlack.PopulationCount())

	for phase := range game_number {
		e.phaseScores[phase] += isolanis[phase]*isolaniDiff +
			doubled[phase]*doublePawnDiff +
			passPawns[phase]*passedPawnDiff +
			backwards[phase]*backwardsPawnDiff +
			opposed[phase]*opposedDiff +
			phalanx[phase]*phalanxDiff +
			supported[phase]*supportedDiff
	}
}
