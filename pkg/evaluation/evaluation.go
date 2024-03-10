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
	maxGamePhase  = 24
	endgameBorder = 12

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
	if pos.HalfMoveClock >= 100 {
		return Contempt(pos)
	}
	if e.evalDraw(pos) {
		return Contempt(pos)
	}
	e.evalPieceSquareTables(pos)
	e.evalKingShield(pos)
	e.evalRooks(pos)
	e.evalPawns(pos)
	e.evalPairs(pos)
	e.evalBaseMaterial(pos)
	e.evalPawnAdjustment(pos)
	e.evalMobilityAndKingAttackValue(pos)
	return e.calculateScore(pos)
}

// Draw Evaluation, https://www.chessprogramming.org/Draw_Evaluation
func (e *eval) evalDraw(pos *position.Position) bool {
	// There are only kings left, it is a draw
	if pos.AllPieces.PopulationCount() == 2 {
		return true
	}

	// If there is any Pawn, Rook or Queen, it is no draw
	if (pos.PiecesBitboard[types.WHITE][types.PAWN] |
		pos.PiecesBitboard[types.BLACK][types.PAWN] |
		pos.PiecesBitboard[types.WHITE][types.ROOK] |
		pos.PiecesBitboard[types.BLACK][types.ROOK] |
		pos.PiecesBitboard[types.WHITE][types.QUEEN] |
		pos.PiecesBitboard[types.BLACK][types.QUEEN]).PopulationCount() > 0 {
		return false
	}
	// At this point there are now Pawns, Rooks or Queens on the board

	numberOfPieces := [types.COLOR_NUMBER]int{
		pos.AllPiecesByColor[types.WHITE].PopulationCount(),
		pos.AllPiecesByColor[types.BLACK].PopulationCount(),
	}

	// If both side have only one minor piece
	if numberOfPieces[types.WHITE] == 2 && numberOfPieces[types.BLACK] == 2 {
		return true
	}

	// If both side have more than one minor piece
	if numberOfPieces[types.WHITE] > 2 && numberOfPieces[types.BLACK] > 2 {
		return false
	}

	// One Side has more than 2 minor pieces
	if numberOfPieces[types.WHITE] > 3 || numberOfPieces[types.BLACK] > 3 {
		return false
	}

	// At this point one side has 2 minor pieces and the other one has only one minor piece.
	// Everything is now a draw, except two Bishops against a Knight.
	for color := types.WHITE; color < types.COLOR_NUMBER; color++ {
		we := color
		them := types.SwitchColor(we)

		// To Bishops agains another minor piece is a draw except
		if pos.PiecesBitboard[we][types.BISHOP].PopulationCount() == 2 {
			return pos.PiecesBitboard[them][types.BISHOP].PopulationCount() == 1
		}
	}

	return true

}

// evalscore calculates the actual score based on the base score and the scores for the different game phases.
func (e *eval) calculateScore(pos *position.Position) int {
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

func gamePhase(pos *position.Position) int {
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
	return gamePhase
}

func IsEndgame(pos *position.Position) bool {
	return gamePhase(pos) < endgameBorder
}
