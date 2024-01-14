package perft

import (
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

type Result struct {
	Move     move.Move
	Position *position.Position
	Leafs    int
}

func Perft(pos *position.Position, depth int) int {
	if depth == 0 {
		return 1
	}
	var nodes int
	moves := pos.GeneratePseudoLegalMoves()
	for _, m := range moves {
		newPos := pos.MakeMove(m)
		if newPos.IsLegal() {
			nodes += Perft(newPos, depth-1)
		}
	}
	return nodes
}

func Divided(pos *position.Position, depth int) []Result {
	moves := pos.GeneratePseudoLegalMoves()
	if depth == 0 {
		return nil
	}
	results := make([]Result, 0, len(moves))
	for _, m := range moves {
		newPos := pos.MakeMove(m)
		if newPos.IsLegal() {
			results = append(results,
				Result{
					Move:     m,
					Position: newPos,
					Leafs:    Perft(newPos, depth-1),
				},
			)
		}
	}
	return results
}
