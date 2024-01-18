package perft

import (
	"fmt"
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func benchmark(b *testing.B, pos *position.Position, depth int) int {
	return Perft(pos, depth)
}

func BenchmarkPerftInitialPosition(b *testing.B) { benchmark(b, position.New(), 5) }

func BenchmarkPerftKiwipete(b *testing.B) {
	pos, err := position.NewFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	assert.NoError(b, err)
	benchmark(b, pos, 5)
}

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
		{
			name:     "initial position",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    5,
			expected: 4865609,
		},
		// {
		// 	name:     "initial position",
		// 	fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// 	depth:    6,
		// 	expected: 119060324,
		// },
		{
			name:     "kiwipete",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    1,
			expected: 48,
		},
		{
			name:     "kiwipete",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    2,
			expected: 2039,
		},
		{
			name:     "kiwipete",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    3,
			expected: 97862,
		},
		{
			name:     "kiwipete",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    4,
			expected: 4085603,
		},
		{
			name:     "kiwipete",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    5,
			expected: 193690690,
		},
		// {
		// 	name:     "kiwipete",
		// 	fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		// 	depth:    6,
		// 	expected: 8031647685,
		// },
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%v", tt.name, tt.depth)
		t.Run(name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, Perft(pos, tt.depth))
		})
	}
}
