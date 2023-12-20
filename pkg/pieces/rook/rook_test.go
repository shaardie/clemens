package rook

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
				occupied: bitboard.BitBySquares(bitboard.SQUARE_C4, bitboard.SQUARE_D5),
			},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_C4,
				bitboard.SQUARE_D1,
				bitboard.SQUARE_D2,
				bitboard.SQUARE_D3,
				bitboard.SQUARE_D5,
				bitboard.SQUARE_E4,
				bitboard.SQUARE_F4,
				bitboard.SQUARE_G4,
				bitboard.SQUARE_H4,
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
