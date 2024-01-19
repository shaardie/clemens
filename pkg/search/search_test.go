package search

import (
	"fmt"
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		depth    int
		expected string
	}{
		{
			name:     "position4",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    5,
			expected: "c4c5",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%v", tt.name, tt.depth)
		t.Run(name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			assert.Equal(t, "c4c5", fmt.Sprintf("%v", Search(pos, tt.depth)))
		})
	}
}
