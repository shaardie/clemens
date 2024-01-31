package search

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/pvline"
	"github.com/shaardie/clemens/pkg/search/transpositiontable"
)

const (
	inf = 1e5
)

type Search struct {
	pos   position.Position
	score int
	nodes uint64
	m     *sync.Mutex
	PV    pvline.PVLine
}

type Info struct {
	Depth uint8
	Time  int64
	Nodes uint64
	PV    pvline.PVLine
	Score int
}

func (s *Search) BestMove() move.Move {
	s.m.Lock()
	defer s.m.Unlock()
	return s.PV.GetBestMove()
}

func NewSearch(pos position.Position) *Search {
	return &Search{
		pos: pos,
		m:   &sync.Mutex{},
	}
}

func (s *Search) Search(ctx context.Context, maxDepth uint8, info chan Info) {
	s.SearchIterative(ctx, maxDepth, info)
}

func (s *Search) SearchIterative(ctx context.Context, maxDepth uint8, info chan Info) {
	for depth := uint8(1); depth <= maxDepth; depth++ {
		i := s.SearchRoot(ctx, depth, -inf, inf)

		select {
		case info <- i:
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (s *Search) SearchRoot(ctx context.Context, depth uint8, alpha, beta int) Info {
	start := time.Now()
	pos := s.pos
	pvl := pvline.PVLine{}
	score := s.negamax(ctx, &pos, -math.MaxInt, math.MaxInt, depth, true, &pvl, true)
	s.PV = pvl
	return Info{
		Depth: depth,
		Time:  time.Since(start).Milliseconds(),
		Nodes: s.nodes,
		Score: score,
		PV:    *pvl.Copy(),
	}
}

func (s *Search) negamax(ctx context.Context, pos *position.Position, alpha, beta int, depth uint8, pvNode bool, pvl *pvline.PVLine, isRoot bool) int {
	s.nodes++

	// Evaluate the leaf node
	if depth == 0 {
		return pos.Evaluation()
	}

	// Check if we can use the transition table but not on root
	if !isRoot {
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
	}
	isRoot = false

	oldAlpha := alpha
	potentialPVLine := pvline.PVLine{}
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
			score := s.negamax(ctx, pos, -beta, -alpha, depth-1, pvNode, &potentialPVLine, isRoot)
			score *= -1
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
				pvl.Update(bestMove, &potentialPVLine)
			}
		}
		*pos = prevPos
		potentialPVLine.Reset()
	}

	nt := transpositiontable.PVNode
	if oldAlpha == alpha {
		nt = transpositiontable.AlphaNode
	}
	transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, depth, alpha, nt)
	return alpha
}
