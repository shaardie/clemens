package perft

import (
	"fmt"
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func TestPerft(t *testing.T) {

	tests := []struct {
		name     string
		fen      string
		depth    int
		expected int
	}{
		{
			name:     "initial position",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    1,
			expected: 20,
		},
		{
			name:     "initial position",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    2,
			expected: 400,
		},
		{
			name:     "initial position",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    3,
			expected: 8902,
		},
		{
			name:     "initial position",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    4,
			expected: 197281,
		},
		// TODO: make them work
		// {
		// 	name:     "initial position",
		// 	fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// 	depth:    5,
		// 	expected: 4865609,
		// },
		// {
		// 	name:     "kiwipete",
		// 	fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		// 	depth:    1,
		// 	expected: 48,
		// },
	}
	for _, tt := range tests {
		name := fmt.Sprintf("perft%v", tt.depth)
		t.Run(name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, Perft(pos, tt.depth))
		})
	}
}
