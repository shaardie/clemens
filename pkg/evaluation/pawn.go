package evaluation

import (
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	midgameIsolanis    = -20
	endgameIsolanis    = -5
	midgameDoubledPawn = -5
	endgameDoubledPawn = -15
)

// evalPawns evaluates the pawn structure
func (e *eval) evalPawns(pos *position.Position) {
	whitePawns := pos.PiecesBitboard[types.WHITE][types.PAWN]
	blackPawns := pos.PiecesBitboard[types.BLACK][types.PAWN]

	// Isolanis
	isolaniDiff := pawn.NumberOfIsolanis(whitePawns) - pawn.NumberOfIsolanis(blackPawns)
	e.phaseScores[midgame] += int16(midgameIsolanis * isolaniDiff)
	e.phaseScores[endgame] += int16(endgameIsolanis * isolaniDiff)

	// Double Pawns
	doublePawnDiff := pawn.NumberOfDoubledPawns(whitePawns) - pawn.NumberOfDoubledPawns(blackPawns)
	e.phaseScores[midgame] += int16(midgameDoubledPawn * doublePawnDiff)
	e.phaseScores[endgame] += int16(endgameDoubledPawn * doublePawnDiff)
}
