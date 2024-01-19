package search

import (
	"math"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

func negaMax(pos *position.Position, depth int) int {
	if depth == 0 {
		return pos.Evaluation()
	}

	max := math.MinInt
	moves := pos.GeneratePseudoLegalMoves()
	var prevPos position.Position
	for _, m := range moves {
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score := -negaMax(pos, depth-1)
			if score > max {
				max = score
			}
		}
		*pos = prevPos
	}
	return max
}

func Search(pos *position.Position, depth int) move.Move {
	if depth == 0 {
		panic("depth should be bigger than 0")
	}
	moves := pos.GeneratePseudoLegalMoves()
	var result move.Move
	max := math.MinInt
	for _, m := range moves {
		prevPos := *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score := -negaMax(pos, depth-1)
			if score > max {
				max = score
				result = m
			}
		}
		pos = &prevPos
	}
	return result
}
