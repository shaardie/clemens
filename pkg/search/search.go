package search

import (
	"math"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

func search(pos *position.Position, alpha, beta, depth int) int {
	if depth == 0 {
		return pos.Evaluation()
	}

	moves := pos.GeneratePseudoLegalMoves()
	var prevPos position.Position
	for _, m := range moves {
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score := -search(pos, -beta, -alpha, depth-1)
			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
		*pos = prevPos
	}
	return alpha
}

func Search(pos *position.Position, depth int) move.Move {
	if depth == 0 {
		panic("depth should be bigger than 0")
	}
	moves := pos.GeneratePseudoLegalMoves()
	var result move.Move
	max := -math.MaxInt
	for _, m := range moves {
		prevPos := *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score := -search(pos, -math.MaxInt, math.MaxInt, depth-1)
			if score >= max {
				max = score
				result = m
			}
		}
		pos = &prevPos
	}
	return result
}
