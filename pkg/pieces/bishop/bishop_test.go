package bishop

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square   uint8
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
				square:   types.SQUARE_D4,
				occupied: bitboard.BitBySquares(types.SQUARE_E5, types.SQUARE_C3),
			},
			want: bitboard.BitBySquares(
				types.SQUARE_A7,
				types.SQUARE_B6,
				types.SQUARE_C3,
				types.SQUARE_C5,
				types.SQUARE_E3,
				types.SQUARE_E5,
				types.SQUARE_F2,
				types.SQUARE_G1,
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
