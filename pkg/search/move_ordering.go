package search

import (
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	pvMoveScore     = 1000
	ttMoveScore     = 900
	killerMoveScore = 100
	promotionScore  = 500
	castlingScore   = 10000
)

// Static Values for MVV-LVA Ordering
// See https://www.chessprogramming.org/MVV-LVA
var MVV_LVA_SCORES [types.PIECE_TYPE_NUMBER - 1][types.PIECE_TYPE_NUMBER]uint16

func init() {
	// Init the MVV-LVA Values
	// For the values to be disjunct, the victim is multiplied by 10
	// To make a difference for PAWNs (value 0) victim a increased by 1
	victim := types.QUEEN
	for {
		for aggressor := types.PAWN; aggressor < types.PIECE_TYPE_NUMBER; aggressor++ {
			MVV_LVA_SCORES[victim][aggressor] = uint16(
				(10*(victim+1) - (aggressor)) + killerMoveScore,
			)
		}
		if victim == types.PAWN {
			break
		}
		victim--
	}
}

func (s *Search) orderMoves(pos *position.Position, moves *move.MoveList, pvMove, ttMove move.Move, ply uint8) {
	for idx := uint8(0); idx < moves.Length(); idx++ {
		m := moves.Get(idx)

		// Use the Move from the principal variation first
		if *m == pvMove {
			m.SetScore(pvMoveScore)
			continue
		}

		// Use the Move from the transposition second
		if *m == ttMove {
			m.SetScore(ttMoveScore)
			continue
		}

		// No capture
		target := pos.GetPiece(m.GetTargetSquare())
		if target == types.NO_PIECE {
			if m.GetMoveType() == move.PROMOTION {
				m.SetScore(promotionScore)
				continue
			}

			// // Killer moves
			if *m == s.KillerMoves[ply][0] {
				m.SetScore(killerMoveScore)
				continue
			}
			if *m == s.KillerMoves[ply][1] {
				m.SetScore(killerMoveScore - 1)
				continue
			}

			// Castling is a little bit good
			if m.GetMoveType() == move.CASTLING {
				m.SetScore(castlingScore)
				continue
			}
			continue
		}

		// Captures are placed above the killer moves. So add the killer move score to all of them.
		source := pos.GetPiece(m.GetSourceSquare())
		m.SetScore(MVV_LVA_SCORES[target.Type()][source.Type()] + promotionScore)
	}
	moves.Sort()
}
