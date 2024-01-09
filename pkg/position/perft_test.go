package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func perft(t *testing.T, pos *Position, depth int) int {
	if depth == 0 {
		return 1
	}

	var nodes int
	moves := pos.GeneratePseudoLegalMoves()
	for _, m := range moves {
		newPos := pos.MakeMove(m)
		if !newPos.IsCheck(pos.sideToMove) {
			t.Log(newPos.ToFen())
			nodes += perft(t, newPos, depth-1)
		}
	}
	return nodes
}

func TestPerft(t *testing.T) {
	assert.Equal(t, 1, perft(t, New(), 0))
	assert.Equal(t, 2, perft(t, New(), 1))
}
