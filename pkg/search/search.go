package search

import (
	"math"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/transpositiontable"
)

type SearchResult struct {
	Score int
	Move  move.Move
}

func search(pos *position.Position, alpha, beta int, depth uint8, pvNode bool) int {
	// Evaluate the leaf node
	if depth == 0 {
		return pos.Evaluation()
	}

	// Check if we can use the transition table
	te, found := transpositiontable.TTable.Get(pos.ZobristHash, depth)
	if found {
		switch te.NodeType {
		case transpositiontable.AlphaNode:
			// return the bigger value of alpha and score
			if te.Score < alpha {
				return alpha
			}
			return te.Score
		case transpositiontable.BetaNode:
			// return the smaller value of beta and score
			if te.Score > beta {
				return beta
			}
			return te.Score
		case transpositiontable.PVNode:
			// return exact value
			return te.Score
		}
	}

	oldAlpha := alpha

	// Generate all moves
	moves := pos.GeneratePseudoLegalMoves()
	var prevPos position.Position

	var bestMove move.Move

	for _, m := range moves {
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score := -search(pos, -beta, -alpha, depth-1, pvNode)
			if pvNode {
				pvNode = false
			}

			if score >= beta {
				transpositiontable.TTable.PotentiallySave(prevPos.ZobristHash, bestMove, depth, beta, transpositiontable.BetaNode)
				return beta
			}

			if score > alpha {
				alpha = score
				bestMove = m
			}
		}
		*pos = prevPos
	}

	nt := transpositiontable.PVNode
	if oldAlpha == alpha {
		nt = transpositiontable.AlphaNode
	}
	transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, depth, alpha, nt)
	return alpha
}

func Search(pos *position.Position, depth uint8) SearchResult {
	if depth == 0 {
		panic("depth should be bigger than 0")
	}
	var currentDepth uint8 = 1
	var score int
	pvNode := true
	r := SearchResult{}
	for {
		r.Score = -math.MaxInt
		moves := pos.GeneratePseudoLegalMoves()
		for _, m := range moves {
			prevPos := *pos
			pos.MakeMove(m)
			if pos.IsLegal() {
				score = -search(pos, -math.MaxInt, math.MaxInt, currentDepth-1, pvNode)
				if score >= r.Score {
					r.Score = score
					r.Move = m
				}
			}
			pos = &prevPos
			if pvNode {
				pvNode = false
			}
		}
		if currentDepth == depth {
			break
		}
		currentDepth++
	}
	return r
}
