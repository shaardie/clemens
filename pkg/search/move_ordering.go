package search

import (
	"math"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const killerMoveScore = 100

// Static Values for MVV-LVA Ordering
// See https://www.chessprogramming.org/MVV-LVA
var MVV_LVA_SCORES [types.PIECE_TYPE_NUMBER - 1][types.PIECE_TYPE_NUMBER]uint16

func init() {
	// Init the MVV-LVA Values
	// For the values to be disjunct, the victim is multiplied by 10
	// To make a difference for PAWNs (value 0) victim a increased by 1
	for victim := types.QUEEN; victim >= types.PAWN; victim-- {
		for aggressor := types.PAWN; aggressor < types.PIECE_TYPE_NUMBER; aggressor++ {
			MVV_LVA_SCORES[victim][aggressor] = uint16(
				(10*(victim+1) - (aggressor)) + killerMoveScore,
			)
		}
	}
}

func MVVLVASort(pos *position.Position, moves *move.MoveList, principalVariationMove, transpositionTableMove move.Move) {
	for idx := uint8(0); idx < moves.Length(); idx++ {
		m := moves.Get(idx)
		if *m == principalVariationMove {
			m.SetScore(math.MaxUint16)
			continue
		}
		if *m == transpositionTableMove {
			m.SetScore(math.MaxUint16 - 1)
			continue
		}

		// Score captures
		target := pos.GetPiece(m.GetTargetSquare())
		// No capture no score
		if target == types.NO_PIECE {
			m.SetScore(0)
			continue
		}
		source := pos.GetPiece(m.GetSourceSquare())
		m.SetScore(MVV_LVA_SCORES[target.Type()][source.Type()])
	}

	moves.Sort()
}

func (s *Search) orderMoves(pos *position.Position, moves *move.MoveList, principalVariationMove, transpositionTableMove move.Move, ply uint8) {
	for idx := uint8(0); idx < moves.Length(); idx++ {
		m := moves.Get(idx)

		// Use the Move from the principal variation first
		if *m == principalVariationMove {
			m.SetScore(math.MaxUint16)
			continue
		}

		// Use the Move from the transposition second
		if *m == transpositionTableMove {
			m.SetScore(math.MaxUint16 - 1)
			continue
		}

		// No capture
		target := pos.GetPiece(m.GetTargetSquare())
		if target == types.NO_PIECE {
			// Killer moves
			if *m == s.KillerMoves[ply][0] {
				m.SetScore(killerMoveScore)
			} else if *m == s.KillerMoves[ply][1] {
				m.SetScore(killerMoveScore - 1)
			}
			continue
		}

		// Captures are placed above the killer moves. So add the killer move score to all of them.
		source := pos.GetPiece(m.GetSourceSquare())
		m.SetScore(MVV_LVA_SCORES[target.Type()][source.Type()] + killerMoveScore)
	}
	moves.Sort()
}
