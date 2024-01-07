package position

import (
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPosition_MakeMove(t *testing.T) {
	tests := []struct {
		name      string
		beforeFen string
		m         move.Move
		afterFen  string
	}{
		{
			name:      "pawn double push",
			beforeFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E2).SetDestinationSquare(types.SQUARE_E4),
			afterFen:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name:      "castling white queen side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E1).SetDestinationSquare(types.SQUARE_C1).SetMoveType(move.CASTLING),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/2KR3R b kq - 0 1",
		},
		{
			name:      "castling white king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E1).SetDestinationSquare(types.SQUARE_G1).SetMoveType(move.CASTLING),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R4RK1 b kq - 0 1",
		},
		{
			name:      "castling black queen side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E8).SetDestinationSquare(types.SQUARE_C8).SetMoveType(move.CASTLING),
			afterFen:  "2kr3r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQ - 0 2",
		},
		{
			name:      "castling black king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E8).SetDestinationSquare(types.SQUARE_G8).SetMoveType(move.CASTLING),
			afterFen:  "r4rk1/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQ - 0 2",
		},
		{
			name:      "update castling rights on white rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_H1).SetDestinationSquare(types.SQUARE_G1),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K1R1 b Qkq - 0 1",
		},
		{
			name:      "update castling rights on white rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_A1).SetDestinationSquare(types.SQUARE_B1),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/1R2K2R b Kkq - 0 1",
		},
		{
			name:      "update castling rights on black rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_H8).SetDestinationSquare(types.SQUARE_G8),
			afterFen:  "r3k1r1/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQq - 0 2",
		},
		{
			name:      "update castling rights on white rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_A8).SetDestinationSquare(types.SQUARE_B8),
			afterFen:  "1r2k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQk - 0 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.beforeFen)
			assert.NoError(t, err)

			assert.Equal(t, tt.afterFen, pos.MakeMove(tt.m).ToFen())
		})
	}
}
