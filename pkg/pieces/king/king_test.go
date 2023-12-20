package king

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
				bitboard.SQUARE_C3,
				bitboard.SQUARE_C4,
				bitboard.SQUARE_C5,
				bitboard.SQUARE_D3,
				bitboard.SQUARE_D5,
				bitboard.SQUARE_E3,
				bitboard.SQUARE_E4,
				bitboard.SQUARE_E5,
			),
		},
		{
			name: "edge",
			args: args{
				square: bitboard.SQUARE_A1,
			},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_A2,
				bitboard.SQUARE_B1,
				bitboard.SQUARE_B2,
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
