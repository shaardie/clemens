package king

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square uint8
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "middle",
			args: args{
				square: types.SQUARE_D4,
			},
			want: bitboard.BitBySquares(
				types.SQUARE_C3,
				types.SQUARE_C4,
				types.SQUARE_C5,
				types.SQUARE_D3,
				types.SQUARE_D5,
				types.SQUARE_E3,
				types.SQUARE_E4,
				types.SQUARE_E5,
			),
		},
		{
			name: "edge",
			args: args{
				square: types.SQUARE_A1,
			},
			want: bitboard.BitBySquares(
				types.SQUARE_A2,
				types.SQUARE_B1,
				types.SQUARE_B2,
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
