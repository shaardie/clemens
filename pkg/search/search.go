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

type Search struct {
	evalMoveList *move.MoveList
}

func NewSearch() *Search {
	return &Search{
		evalMoveList: move.NewMoveList(),
	}
}

func (s *Search) search(pos *position.Position, alpha, beta int, depth uint8, pvNode bool) int {
	// Evaluate the leaf node
	if depth == 0 {
		s.evalMoveList.Reset()
		return pos.Evaluation(s.evalMoveList)
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

	var prevPos position.Position
	var bestMove move.Move

	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score := -s.search(pos, -beta, -alpha, depth-1, pvNode)
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

func (s *Search) Search(pos *position.Position, depth uint8) SearchResult {
	if depth == 0 {
		panic("depth should be bigger than 0")
	}
	var currentDepth uint8 = 1
	var score int
	pvNode := true
	r := SearchResult{}
	for {
		r.Score = -math.MaxInt
		// Generate all moves
		moves := move.NewMoveList()
		pos.GeneratePseudoLegalMoves(moves)
		for i := uint8(0); i < moves.Length(); i++ {
			m := moves.Get(i)
			prevPos := *pos
			pos.MakeMove(m)
			if pos.IsLegal() {
				score = -s.search(pos, -math.MaxInt, math.MaxInt, currentDepth-1, pvNode)
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
