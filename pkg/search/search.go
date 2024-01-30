package search

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/transpositiontable"
)

type Search struct {
	pos      position.Position
	score    int
	bestMove move.Move
	nodes    uint64
	pvNodes  []move.Move
	m        *sync.Mutex
}

type Info struct {
	Depth uint8
	Time  int64
	Nodes uint64
	PV    []move.Move
	Score int
}

func (s *Search) BestMove() move.Move {
	s.m.Lock()
	defer s.m.Unlock()
	return s.bestMove
}

func NewSearch(pos position.Position) *Search {
	return &Search{
		pos: pos,
		m:   &sync.Mutex{},
	}
}

func (s *Search) Search(ctx context.Context, maxDepth uint8, info chan Info) {
	var currentDepth uint8 = 0
	var isPVNode bool
	start := time.Now()
	var pos position.Position
	pos = s.pos
	s.score = -math.MaxInt

	for {

		// Stop on abort
		select {
		case <-ctx.Done():
			return
		default:
		}

		if maxDepth > 0 && currentDepth == maxDepth {
			return
		}
		currentDepth++
		isPVNode = true
		s.pvNodes = make([]move.Move, 0, currentDepth)

		// Generate all moves
		moves := move.NewMoveList()
		pos.GeneratePseudoLegalMoves(moves)
		for i := uint8(0); i < moves.Length(); i++ {
			m := moves.Get(i)
			prevPos := pos
			pos.MakeMove(m)
			if pos.IsLegal() {
				if isPVNode {
					s.pvNodes = append(s.pvNodes, m)
				}

				s.nodes++
				score := -s.search(ctx, &pos, -math.MaxInt, math.MaxInt, currentDepth-1, isPVNode)
				if score > s.score {
					s.score = score
					s.bestMove = m
				}
				if isPVNode {
					isPVNode = false
				}

			}
			pos = prevPos
		}

		if info != nil {
			select {
			case info <- Info{
				Depth: currentDepth,
				Time:  time.Since(start).Milliseconds(),
				Nodes: s.nodes,
				Score: s.score,
			}:
			default:
				panic("info channel broken")
			}
		}
	}
}
func (s *Search) search(ctx context.Context, pos *position.Position, alpha, beta int, depth uint8, pvNode bool) int {
	// Stop on abort
	select {
	case <-ctx.Done():
		// not quite sure is alpha is a good value here
		return alpha
	default:
	}

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
			s.nodes++
			score := -s.search(ctx, pos, -beta, -alpha, depth-1, pvNode)
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
