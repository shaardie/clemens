package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	midgame int = iota
	endgame
	game_number
)

const (
	maxGamePhase         = 24
	knightGamePhaseValue = 1
	bishopGamePhaseValue = 1
	RookGamePhaseValue   = 2
	QueenGamePhaseValue  = 4
)

// Evaluation evaluates the position.
func Evaluation(pos *position.Position) int {
	return evalWithCache(pos)
}

// evalWithCache uses a transposition table to lookup, if a position was already evaluted.
// If so, it returns the cached value other calls the actual evaluation and safe the result in the cache.
func evalWithCache(pos *position.Position) int {
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
	phaseScores [game_number]int
	baseScore   int
}

// do is the actual evaluation function
func (e *eval) do(pos *position.Position) int {
	e.evalPieceSquareTables(pos)
	e.evalKingShield(pos)
	e.evalRooks(pos)
	e.evalBishop(pos)
	e.evalPawns(pos)
	e.evalPairs(pos)
	e.evalBaseMaterial(pos)
	e.evalPawnAdjustment(pos)
	e.evalMobilityAndKingAttackValue(pos)
	return e.calculateScore(pos)
}

// evalscore calculates the actual score based on the base score and the scores for the different game phases.
func (e *eval) calculateScore(pos *position.Position) int {
	// Calculate the game phase based on the number of specific PieceTypes maxed by maxGamePhase.
	var gamePhase int
	for color := types.WHITE; color < types.COLOR_NUMBER; color++ {
		gamePhase += bishopGamePhaseValue * pos.PiecesBitboard[color][types.BISHOP].PopulationCount()
		gamePhase += knightGamePhaseValue * pos.PiecesBitboard[color][types.KNIGHT].PopulationCount()
		gamePhase += RookGamePhaseValue * pos.PiecesBitboard[color][types.ROOK].PopulationCount()
		gamePhase += QueenGamePhaseValue * pos.PiecesBitboard[color][types.QUEEN].PopulationCount()
	}
	if gamePhase > maxGamePhase {
		gamePhase = maxGamePhase
	}

	// Merge midgame and endgame value proportionally and add base score
	score := (e.phaseScores[midgame]*gamePhase + e.phaseScores[endgame]*(maxGamePhase-gamePhase)) / maxGamePhase
	score += e.baseScore

	// Make the result side aware
	if pos.SideToMove == types.BLACK {
		score *= -1
	}
	return score
}
