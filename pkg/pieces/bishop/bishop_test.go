package bishop

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square   int
		occupied bitboard.Bitboard
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "middle",
			args: args{
				square:   bitboard.SQUARE_D4,
				occupied: bitboard.BitBySquares(bitboard.SQUARE_E5, bitboard.SQUARE_C3),
			},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_A7,
				bitboard.SQUARE_B6,
				bitboard.SQUARE_C3,
				bitboard.SQUARE_C5,
				bitboard.SQUARE_E3,
				bitboard.SQUARE_E5,
				bitboard.SQUARE_F2,
				bitboard.SQUARE_G1,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttacksBySquare(tt.args.square, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttacksBySquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
