package pieces

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
)

func TestSlidingAttacks(t *testing.T) {
	type args struct {
		t        pieceType
		square   int
		occupied bitboard.Bitboard
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "rook",
			args: args{t: RookType, square: bitboard.SQUARE_A1, occupied: bitboard.RankMask6 | bitboard.FileMaskE},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_A2,
				bitboard.SQUARE_A3,
				bitboard.SQUARE_A4,
				bitboard.SQUARE_A5,
				bitboard.SQUARE_A6,
				bitboard.SQUARE_B1,
				bitboard.SQUARE_C1,
				bitboard.SQUARE_D1,
				bitboard.SQUARE_E1,
			),
		},
		{
			name: "bishop",
			args: args{t: BishopType, square: bitboard.SQUARE_B2, occupied: bitboard.RankMask6 | bitboard.FileMaskE},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_A1,
				bitboard.SQUARE_A3,
				bitboard.SQUARE_C1,
				bitboard.SQUARE_C3,
				bitboard.SQUARE_D4,
				bitboard.SQUARE_E5,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlidingAttacks(tt.args.t, tt.args.square, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RookAttacks() = %v, want %v", got.PrettyString(), tt.want.PrettyString())
			}
		})
	}
}
