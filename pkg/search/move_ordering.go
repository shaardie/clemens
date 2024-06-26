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
	couterMoveBonus = 10
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

func (s *Search) scoreMoves(pos *position.Position, moves *move.MoveList, pvMove, ttMove move.Move, ply uint8) {
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
		targetSquare := m.GetTargetSquare()
		sourceSquare := m.GetSourceSquare()
		var target types.Piece
		if m.GetMoveType() == move.EN_PASSANT {
			if pos.SideToMove == types.BLACK {
				target = types.WHITE_PAWN
			} else {
				target = types.BLACK_PAWN
			}
		} else {
			target = pos.GetPiece(targetSquare)
		}

		if target == types.NO_PIECE {
			if m.GetMoveType() == move.PROMOTION {
				m.SetScore(promotionScore)
				continue
			}

			// Killer moves
			if *m == s.KillerMoves[ply][0] {
				m.SetScore(killerMoveScore)
				continue
			}
			if *m == s.KillerMoves[ply][1] {
				m.SetScore(killerMoveScore - 1)
				continue
			}

			score := s.history[pos.SideToMove][sourceSquare][targetSquare]
			counterMove := s.counter[pos.SideToMove][sourceSquare][targetSquare]
			if counterMove == *m {
				score += couterMoveBonus
			}

			m.SetScore(score)
			continue

		}

		// Captures are placed above the killer moves. So add the killer move score to all of them.
		source := pos.GetPiece(sourceSquare)
		m.SetScore(MVV_LVA_SCORES[target.Type()][source.Type()] + promotionScore)
	}
}

func (s *Search) orderMoves(pos *position.Position, moves *move.MoveList, pvMove, ttMove move.Move, ply uint8) {
	s.scoreMoves(pos, moves, pvMove, ttMove, ply)
	moves.Sort()
}
