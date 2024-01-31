package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSearchKiwipete(b *testing.B) {
	pos, err := position.NewFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	assert.NoError(b, err)

	NewSearch(*pos).Search(context.TODO(), 6, nil)
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name        string
		fen         string
		depth       uint8
		notExpected string
		expected    string
	}{
		{
			name:        "startpos",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:       5,
			notExpected: "a1a1",
		},
		{
			name:        "position4",
			fen:         "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:       5,
			notExpected: "a1a1",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%v", tt.name, tt.depth)
		t.Run(name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			s := NewSearch(*pos)
			s.Search(context.TODO(), tt.depth, nil)
			if tt.notExpected != "" {
				assert.NotEqual(t, tt.notExpected, fmt.Sprintf("%v", s.m))
			}
			if tt.expected != "" {
				assert.Equal(t, tt.expected, fmt.Sprintf("%v", s.m))
			}
		})
	}
}
