package knight

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square int
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "middle",
			args: args{
				square: bitboard.SQUARE_D4,
			},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_B3,
				bitboard.SQUARE_B5,
				bitboard.SQUARE_C2,
				bitboard.SQUARE_C6,
				bitboard.SQUARE_E2,
				bitboard.SQUARE_E6,
				bitboard.SQUARE_F3,
				bitboard.SQUARE_F5,
			),
		},
		{
			name: "edge",
			args: args{
				square: bitboard.SQUARE_B1,
			},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_A3,
				bitboard.SQUARE_C3,
				bitboard.SQUARE_D2,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttacksBySquare(tt.args.square); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttacksBySquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
