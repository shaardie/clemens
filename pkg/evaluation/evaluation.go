package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	midgame int = iota
	endgame
	game_number
)

const (
	knightGamePhaseValue = 1
	bishopGamePhaseValue = 1
	rookGamePhaseValue   = 2
	queenGamePhaseValue  = 4

	maxGamePhase  = 4*knightGamePhaseValue + 4*bishopGamePhaseValue + 4*rookGamePhaseValue + 2*queenGamePhaseValue
	endgameBorder = maxGamePhase / 2
)

// Evaluation evaluates the position.
func Evaluation(pos *position.Position) int16 {
	return evalWithCache(pos)
}

// evalWithCache uses a transposition table to lookup, if a position was already evaluted.
// If so, it returns the cached value other calls the actual evaluation and safe the result in the cache.
func evalWithCache(pos *position.Position) int16 {
	score, found := tTable.get(pos.ZobristHash)
	if found {
		return score
	}

	e := eval{}
	score = e.do(pos)

	// Save score in transposition table
	tTable.save(pos.ZobristHash, score)

	return score

}

type eval struct {
	phaseScores [game_number]int16
	baseScore   int16
}

// do is the actual evaluation function
func (e *eval) do(pos *position.Position) int16 {
	if e.isDraw(pos) {
		return Contempt(pos)
	}
	e.evalPieceSquareTables(pos)
	// e.evalKingShield(pos)
	// e.evalRooks(pos)
	e.evalPawns(pos)
	e.evalPairs(pos)
	e.evalBaseMaterial(pos)
	e.evalPawnAdjustment(pos)
	e.evalMobilityAndKingAttackValue(pos)
	return e.calculateScore(pos)
}

// evalscore calculates the actual score based on the base score and the scores for the different game phases.
func (e *eval) calculateScore(pos *position.Position) int16 {
	gamePhase := gamePhase(pos)

	// Merge midgame and endgame value proportionally and add base score
	score := (e.phaseScores[midgame]*gamePhase + e.phaseScores[endgame]*(maxGamePhase-gamePhase)) / maxGamePhase
	score += e.baseScore

	// Make the result side aware
	if pos.SideToMove == types.BLACK {
		score *= -1
	}
	return score
}

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
