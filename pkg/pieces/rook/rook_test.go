package rook

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
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
				square:   SQUARE_D4,
				occupied: bitboard.BitBySquares(SQUARE_C4, SQUARE_D5),
			},
			want: bitboard.BitBySquares(
				SQUARE_C4,
				SQUARE_D1,
				SQUARE_D2,
				SQUARE_D3,
				SQUARE_D5,
				SQUARE_E4,
				SQUARE_F4,
				SQUARE_G4,
				SQUARE_H4,
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
