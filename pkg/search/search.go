package search

import (
	"math"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

type SearchResult struct {
	Score int
	Move  move.Move
	Nodes int
}

func search(pos *position.Position, alpha, beta, depth int) (int, int) {
	if depth == 0 {
		return pos.Evaluation(), 1
	}

	var nodes int
	moves := pos.GeneratePseudoLegalMoves()
	var prevPos position.Position
	for _, m := range moves {
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score, additionalNodes := search(pos, -beta, -alpha, depth-1)
			score = score * -1
			nodes += additionalNodes
			if score >= beta {
				return beta, nodes
			}
			if score > alpha {
				alpha = score
			}
		}
		*pos = prevPos
	}
	return alpha, nodes
}

func Search(pos *position.Position, depth int) SearchResult {
	if depth == 0 {
		panic("depth should be bigger than 0")
	}
	r := SearchResult{
		Score: -math.MaxInt,
	}
	moves := pos.GeneratePseudoLegalMoves()

	for _, m := range moves {
		prevPos := *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score, nodes := search(pos, -math.MaxInt, math.MaxInt, depth-1)
			score = -1 * score
			r.Nodes += nodes
			if score >= r.Score {
				r.Score = score
				r.Move = m
			}
		}
		pos = &prevPos
	}
	return r
}
