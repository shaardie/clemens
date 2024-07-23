package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	knightGamePhaseValue = 1
	bishopGamePhaseValue = 1
	rookGamePhaseValue   = 2
	queenGamePhaseValue  = 4

	maxGamePhase  = 4*knightGamePhaseValue + 4*bishopGamePhaseValue + 4*rookGamePhaseValue + 2*queenGamePhaseValue
	endgameBorder = maxGamePhase / 2
)

func gamePhase(pos *position.Position) int16 {
	// Calculate the game phase based on the number of specific PieceTypes maxed by maxGamePhase.
	var gamePhase int16
	for color := types.WHITE; color < types.COLOR_NUMBER; color++ {
		gamePhase += int16(bishopGamePhaseValue * pos.PiecesBitboard[color][types.BISHOP].PopulationCount())
		gamePhase += int16(knightGamePhaseValue * pos.PiecesBitboard[color][types.KNIGHT].PopulationCount())
		gamePhase += int16(rookGamePhaseValue * pos.PiecesBitboard[color][types.ROOK].PopulationCount())
		gamePhase += int16(queenGamePhaseValue * pos.PiecesBitboard[color][types.QUEEN].PopulationCount())
	}
	if gamePhase > maxGamePhase {
		gamePhase = maxGamePhase
	}
	return gamePhase
}

func IsEndgame(pos *position.Position) bool {
	return gamePhase(pos) < endgameBorder
}

func IsPawnEndgame(pos *position.Position) bool {
	for c := range types.COLOR_NUMBER {
		for t := range types.PIECE_TYPE_NUMBER {
			if pos.PiecesBitboard[c][t] != bitboard.Empty {
				return false
			}
		}
	}
	return true
}
