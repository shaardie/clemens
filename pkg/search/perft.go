package search

import (
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

type PerftResults struct {
	Move     move.Move
	Position position.Position
	Leafs    int
}

func Perft(pos *position.Position, depth int) int {
	if depth == 0 {
		return 1
	}
	var nodes int
	var s position.State
	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		pos.MakeMove(*m, &s)
		if pos.IsLegal() {
			nodes += Perft(pos, depth-1)
		}
		pos.UnMakeMove(&s)
	}
	return nodes
}

func Divided(pos *position.Position, depth int) []PerftResults {
	if depth == 0 {
		panic("depth should be bigger than 0")
	}

	var s position.State

	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	results := make([]PerftResults, 0, moves.Length())
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		pos.MakeMove(*m, &s)
		if pos.IsLegal() {
			results = append(results,
				PerftResults{
					Move:     *m,
					Position: *pos,
					Leafs:    Perft(pos, depth-1),
				},
			)
		}
		pos.UnMakeMove(&s)
	}
	return results
}
