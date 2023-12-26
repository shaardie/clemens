package rook

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
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
				square:   types.SQUARE_D4,
				occupied: bitboard.BitBySquares(types.SQUARE_C4, types.SQUARE_D5),
			},
			want: bitboard.BitBySquares(
				types.SQUARE_C4,
				types.SQUARE_D1,
				types.SQUARE_D2,
				types.SQUARE_D3,
				types.SQUARE_D5,
				types.SQUARE_E4,
				types.SQUARE_F4,
				types.SQUARE_G4,
				types.SQUARE_H4,
			),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttacksBySquare(tt.args.square, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttacksBySquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
