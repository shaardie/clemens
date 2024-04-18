package pawn

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_doublePushTargets(t *testing.T) {
	type args struct {
		c        types.Color
		pawns    bitboard.Bitboard
		occupied bitboard.Bitboard
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "white double push",
			args: args{
				c:        types.WHITE,
				pawns:    bitboard.BitBySquares(types.SQUARE_A2),
				occupied: bitboard.Empty,
			},
			want: bitboard.BitBySquares(types.SQUARE_A4),
		},
		{
			name: "white no double push",
			args: args{
				c:        types.WHITE,
				pawns:    bitboard.BitBySquares(types.SQUARE_A3),
				occupied: bitboard.Empty,
			},
			want: bitboard.Empty,
		},
		{
			name: "blach double push",
			args: args{
				c:        types.BLACK,
				pawns:    bitboard.BitBySquares(types.SQUARE_B7),
				occupied: bitboard.Empty,
			},
			want: bitboard.BitBySquares(types.SQUARE_B5),
		},
		{
			name: "black no double push",
			args: args{
				c:        types.BLACK,
				pawns:    bitboard.BitBySquares(types.SQUARE_B6),
				occupied: bitboard.Empty,
			},
			want: bitboard.Empty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doublePushTargets(tt.args.c, tt.args.pawns, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doublePushTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberOfIsolanis(t *testing.T) {
	tests := []struct {
		pawns bitboard.Bitboard
		want  int
	}{
		{
			pawns: bitboard.BitBySquares(
				types.SQUARE_A2,
				types.SQUARE_C2,
				types.SQUARE_D2,
				types.SQUARE_F2),
			want: 2,
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, NumberOfIsolanis(tt.pawns))
	}
}

func TestNumberDouble(t *testing.T) {
	tests := []struct {
		pawns bitboard.Bitboard
		want  int
	}{
		{
			pawns: bitboard.BitBySquares(
				types.SQUARE_A2, types.SQUARE_A3, types.SQUARE_A4,
				types.SQUARE_C2,
				types.SQUARE_D2, types.SQUARE_D5,
			),
			want: 3,
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, NumberOfIsolanis(tt.pawns))
	}
}

func TestBackwards(t *testing.T) {
	tests := []struct {
		color  types.Color
		wPawns bitboard.Bitboard
		bPawns bitboard.Bitboard
		want   bitboard.Bitboard
	}{
		{
			color:  types.WHITE,
			wPawns: bitboard.BitBySquares(types.SQUARE_C3),
			bPawns: bitboard.BitBySquares(types.SQUARE_D5),
			want:   bitboard.BitBySquares(types.SQUARE_C3),
		}, {
			color:  types.BLACK,
			wPawns: bitboard.BitBySquares(types.SQUARE_C3),
			bPawns: bitboard.BitBySquares(types.SQUARE_D5),
			want:   bitboard.BitBySquares(types.SQUARE_D5),
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, Backwards(tt.color, tt.wPawns, tt.bPawns))
	}
}

func TestBackwardsPawns(t *testing.T) {
	tests := []struct {
		color types.Color
		pawns bitboard.Bitboard
		want  bitboard.Bitboard
	}{
		{
			color: types.WHITE,
			pawns: bitboard.BitBySquares(types.SQUARE_B2, types.SQUARE_C2),
			want:  bitboard.BitBySquares(types.SQUARE_B2, types.SQUARE_C2),
		},
		{
			color: types.WHITE,
			pawns: bitboard.BitBySquares(types.SQUARE_B2),
			want:  bitboard.Empty,
		},
		{
			color: types.WHITE,
			pawns: bitboard.BitBySquares(types.SQUARE_B2, types.SQUARE_C2, types.SQUARE_D2, types.SQUARE_F2, types.SQUARE_H2),
			want:  bitboard.BitBySquares(types.SQUARE_B2, types.SQUARE_C2, types.SQUARE_D2),
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, Phalanx(tt.pawns))
	}
}
