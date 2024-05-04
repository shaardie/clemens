package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_kingShield(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want int16
	}{
		{
			name: "startpos",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: 0,
		},
		{
			name: "white king castle",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 b kq - 0 1",
			want: 3 * shieldValue[shieldRank2],
		},
		{
			name: "black king castle and one pawn up",
			fen:  "2kr1bnr/1ppppppp/p7/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 0 1",
			want: -2*shieldValue[shieldRank2] - shieldValue[shieldRank3],
		},
		{
			name: "white king not on the back rank",
			fen:  "rnbqkbnr/pppppppp/8/8/8/6K1/PPPPPPPP/RNBQ1BNR b KQkq - 0 1",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			require.NoError(t, err)
			assert.Equal(t, tt.want, kingShield(pos))
		})
	}
}
