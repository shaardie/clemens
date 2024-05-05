package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_passed(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want int16
	}{
		{
			name: "one passed white pawn on 4th rank",
			fen:  "rnbqkbnr/pp3ppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			want: 3 * passedScalar,
		},
		{
			name: "one passed black pawn on 5th rank",
			fen:  "rnbqkbnr/ppp1pppp/8/3p4/8/8/PP3PPP/RNBQKBNR w KQkq - 0 1",
			want: -3 * passedScalar,
		},
		{
			name: "doubled passed white pawns on 3th and 4th rank",
			fen:  "rnbqkbnr/pp3ppp/8/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 1",
			want: 3 * passedScalar,
		},
		{
			name: "doubled passed black pawns on 5th and 6th rank",
			fen:  "rnbqkbnr/ppp1pppp/3p4/3p4/8/8/PP3PPP/RNBQKBNR w KQkq - 0 1",
			want: -3 * passedScalar,
		},
		{
			name: "2 passed white pawns on 4th and 5th rank",
			fen:  "rnbqkbnr/pp4pp/8/4P3/3P4/8/PP3PPP/RNBQKBNR w KQkq - 0 1",
			want: 7 * passedScalar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			require.NoError(t, err)
			assert.Equal(t, tt.want, passed(
				pos.PiecesBitboard[types.WHITE][types.PAWN],
				pos.PiecesBitboard[types.BLACK][types.PAWN]))
		})
	}
}

func Test_phalanx(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want int16
	}{
		{
			name: "start posisition",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: 0,
		},
		{
			name: "black no pawns",
			fen:  "rnbqkbnr/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: 8 * phalanxScalar,
		},
		{
			name: "3 black phalanx pawns on the 5th rank",
			fen:  "rnbqkbnr/ppp3pp/8/3ppp2/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: -(3*3 - 3) * phalanxScalar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			require.NoError(t, err)
			assert.Equal(t, tt.want, phalanx(
				pos.PiecesBitboard[types.WHITE][types.PAWN],
				pos.PiecesBitboard[types.BLACK][types.PAWN]))
		})
	}
}

func Test_supported(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want int16
	}{
		{
			name: "start posisition",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: 0,
		},
		{
			name: "supported white pawn in the center",
			fen:  "rnbqkbnr/pppppppp/8/3P4/4P3/2N5/PPP2PPP/R1BQKBNR w KQkq - 0 1",
			want: 4 * supportedScalar,
		},
		{
			name: "3 supported black pawns on 6th, 5th and 4th rank",
			fen:  "rnbqkbnr/ppp3pp/5p2/4p3/3p4/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: -9 * supportedScalar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			require.NoError(t, err)
			assert.Equal(t, tt.want, supported(
				pos.PiecesBitboard[types.WHITE][types.PAWN],
				pos.PiecesBitboard[types.BLACK][types.PAWN]))
		})
	}
}
