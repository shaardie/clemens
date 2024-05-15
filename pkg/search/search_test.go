package search

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkSearchKiwipete(b *testing.B) {
	// Redirect standard out to null
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)

	pos, err := position.NewFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	require.NoError(b, err)
	s := NewSearch(*pos)
	s.Search(context.TODO(), SearchParameter{Depth: 7, Infinite: true})
}

func BenchmarkSearchStartPos(b *testing.B) {
	// Redirect standard out to null
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)

	s := NewSearch(*position.New())
	s.Search(context.TODO(), SearchParameter{Depth: 7, Infinite: true})
}

func TestSearchTimeout(t *testing.T) {
	pos, err := position.NewFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	require.NoError(t, err)
	s := NewSearch(*pos)
	s.Search(context.TODO(), SearchParameter{Depth: 10, MoveTime: 1000})
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
			depth:       2,
			notExpected: "a1a1",
		},
		{
			name:        "position4",
			fen:         "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:       2,
			notExpected: "a1a1",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%v", tt.name, tt.depth)
		t.Run(name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			s := NewSearch(*pos)
			s.Search(context.TODO(), SearchParameter{Depth: tt.depth, Infinite: true})
			if tt.notExpected != "" {
				assert.NotEqual(t, tt.notExpected, s.bestMove().String())
			}
			if tt.expected != "" {
				assert.Equal(t, tt.expected, s.bestMove().String())
			}
		})
	}
}

func Test_calculateTime(t *testing.T) {
	type args struct {
		sideToMove types.Color
		plys       int
		sp         SearchParameter
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "error case found on lichess",
			args: args{
				sideToMove: types.WHITE,
				plys:       0,
				sp: SearchParameter{
					WTime: 60000,
					BTime: 60000,
					WInc:  2000,
					BInc:  2000,
				},
			},
			want: 2700,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTime(tt.args.sideToMove, tt.args.plys, tt.args.sp); got != tt.want {
				t.Errorf("calculateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
