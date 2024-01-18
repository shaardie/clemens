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
		{
			name:      "white doing en passant",
			beforeFen: "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E5).SetDestinationSquare(types.SQUARE_D6).SetMoveType(move.EN_PASSANT),
			afterFen:  "rnbqkbnr/ppp1pppp/3P4/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:      "black doing en passant",
			beforeFen: "rnbqkbnr/pppp1ppp/8/8/3Pp3/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E4).SetDestinationSquare(types.SQUARE_D3).SetMoveType(move.EN_PASSANT),
			afterFen:  "rnbqkbnr/pppp1ppp/8/8/8/3p4/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
		},
		{
			name:      "promote white pawn to queen",
			beforeFen: "8/3k3P/8/8/8/8/8/3K4 w - - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_H7).SetDestinationSquare(types.SQUARE_H8).SetMoveType(move.PROMOTION).SetPromitionPieceType(types.QUEEN),
			afterFen:  "7Q/3k4/8/8/8/8/8/3K4 b - - 0 1",
		},
		{
			name:      "promote black pawn to rook",
			beforeFen: "8/3k4/8/8/8/8/6p1/3K4 b - - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_G2).SetDestinationSquare(types.SQUARE_G1).SetMoveType(move.PROMOTION).SetPromitionPieceType(types.ROOK),
			afterFen:  "8/3k4/8/8/8/8/8/3K2r1 w - - 0 2",
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

func TestPosition_GeneratePseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name      string
		beforeFen string
		moves     []move.Move
	}{{
		name:      "Promotion",
		beforeFen: "r3k2r/p1ppqpb1/bn2pnp1/1N1PN3/1p2P3/5Q2/PPPB1PpP/R3KB1R b KQkq - 0 1",
		moves: []move.Move{
			0xe78, 0xeb8, 0xef8, 0x3ff, 0x5ff, 0x7ff, 0x9ff, 0xbff, 0xdff, 0xf7f,
			0xfbf, 0x868, 0xc68, 0xea8, 0xbf6, 0xf76, 0x8b4, 0xaf4, 0xef4, 0xf74,
			0x629, 0x6a9, 0x8e9, 0xea9, 0x72d, 0x7ad, 0x8ed, 0x9ed, 0xded, 0xfad,
			0xefc, 0xf7c, 0x118e, 0x518e, 0x918e, 0xd18e, 0x114e, 0x514e, 0x914e,
			0xd14e, 0x11ce, 0x51ce, 0x91ce, 0xd1ce, 0x459, 0x8ec, 0x9ae, 0x8b2,
			0xab2, 0xaf3, 0x3fbc, 0x3ebc,
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.beforeFen)
			assert.NoError(t, err)

			assert.Equal(t, tt.moves, pos.GeneratePseudoLegalMoves())
		})
	}
}
