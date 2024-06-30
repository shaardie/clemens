package evaluation

import (
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	isolanis = [game_number]int16{-20, -5}
	passed   = [game_number]int16{10, 25}
)

// evalPawns evaluates the pawn structure
func (e *eval) evalPawns(pos *position.Position) {
	whitePawns := pos.PiecesBitboard[types.WHITE][types.PAWN]
	blackPawns := pos.PiecesBitboard[types.BLACK][types.PAWN]

	isolaniDiff := int16(pawn.NumberOfIsolanis(whitePawns) - pawn.NumberOfIsolanis(blackPawns))
	passedDiff := int16(pawn.NumberOfPassed(types.WHITE, whitePawns, blackPawns) - pawn.NumberOfPassed(types.BLACK, whitePawns, blackPawns))
	for phase := range game_number {
		e.phaseScores[phase] += isolanis[phase]*isolaniDiff + passed[phase]*passedDiff
	}
}
