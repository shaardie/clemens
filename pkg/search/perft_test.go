package search

import (
	"fmt"
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func benchmark(b *testing.B, pos *position.Position, depth int) int {
	return Perft(pos, depth)
}

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
		// {
		// 	name:     "kiwipete",
		// 	fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		// 	depth:    5,
		// 	expected: 193690690,
		// },
		// {
		// 	name:     "kiwipete",
		// 	fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		// 	depth:    6,
		// 	expected: 8031647685,
		// },
		{
			name:     "position3",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    1,
			expected: 14,
		},
		{
			name:     "position3",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    2,
			expected: 191,
		},
		{
			name:     "position3",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    3,
			expected: 2812,
		},
		{
			name:     "position3",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    4,
			expected: 43238,
		},
		{
			name:     "position3",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    5,
			expected: 674624,
		},
		{
			name:     "position4",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    1,
			expected: 6,
		},
		{
			name:     "position4",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    2,
			expected: 264,
		},
		{
			name:     "position4",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    3,
			expected: 9467,
		},
		{
			name:     "position4",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    4,
			expected: 422333,
		},
		{
			name:     "position4",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    5,
			expected: 15833292,
		},
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
