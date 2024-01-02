package position

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPosition_SquareAttackedBy(t *testing.T) {
	type fields struct {
		PiecesBB   [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER]bitboard.Bitboard
		SideToMove types.Color
	}
	type args struct {
		square   int
		occupied bitboard.Bitboard
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bitboard.Bitboard
	}{
		{
			name: "attacked",
			fields: fields{
				PiecesBB: [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER]bitboard.Bitboard{
					{
						bitboard.BitBySquares(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := Position{
				piecesBitboard: tt.fields.PiecesBB,
				sideToMove:     tt.fields.SideToMove,
			}
			if got := pos.SquareAttackedBy(tt.args.square, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.SquareAttackedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_SetPiece(t *testing.T) {
	assert.NoError(t, New().validate())

}
