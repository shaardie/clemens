package position

import (
	"fmt"
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPosition_MakeAndUnMakeNullMove(t *testing.T) {
	pos, err := NewFromFen("rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1")
	require.NoError(t, err)
	ep := pos.EnPassant
	oldPos := *pos
	pos.MakeNullMove()
	newPos := *pos
	pos.UnMakeNullMove(ep)
	assert.Equal(t, oldPos, *pos)
	assert.NotEqual(t, newPos, oldPos)
}

func TestPosition_MakeAndUnmakeMove(t *testing.T) {
	tests := []struct {
		name      string
		beforeFen string
		m         move.Move
		afterFen  string
	}{
		{
			name:      "pawn double push",
			beforeFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E2).SetTargetSquare(types.SQUARE_E4),
			afterFen:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name:      "castling white queen side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E1).SetTargetSquare(types.SQUARE_C1).SetMoveType(move.CASTLING),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/2KR3R b kq - 1 1",
		},
		{
			name:      "castling white king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E1).SetTargetSquare(types.SQUARE_G1).SetMoveType(move.CASTLING),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R4RK1 b kq - 1 1",
		},
		{
			name:      "castling black queen side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E8).SetTargetSquare(types.SQUARE_C8).SetMoveType(move.CASTLING),
			afterFen:  "2kr3r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQ - 1 2",
		},
		{
			name:      "castling black king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E8).SetTargetSquare(types.SQUARE_G8).SetMoveType(move.CASTLING),
			afterFen:  "r4rk1/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQ - 1 2",
		},
		{
			name:      "update castling rights on white rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_H1).SetTargetSquare(types.SQUARE_G1),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K1R1 b Qkq - 1 1",
		},
		{
			name:      "update castling rights on white rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_A1).SetTargetSquare(types.SQUARE_B1),
			afterFen:  "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/1R2K2R b Kkq - 1 1",
		},
		{
			name:      "update castling rights on black rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_H8).SetTargetSquare(types.SQUARE_G8),
			afterFen:  "r3k1r1/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQq - 1 2",
		},
		{
			name:      "update castling rights on white rook move on king side",
			beforeFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_A8).SetTargetSquare(types.SQUARE_B8),
			afterFen:  "1r2k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQk - 1 2",
		},
		{
			name:      "white doing en passant",
			beforeFen: "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E5).SetTargetSquare(types.SQUARE_D6).SetMoveType(move.EN_PASSANT),
			afterFen:  "rnbqkbnr/ppp1pppp/3P4/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:      "black doing en passant",
			beforeFen: "rnbqkbnr/pppp1ppp/8/8/3Pp3/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_E4).SetTargetSquare(types.SQUARE_D3).SetMoveType(move.EN_PASSANT),
			afterFen:  "rnbqkbnr/pppp1ppp/8/8/8/3p4/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
		},
		{
			name:      "promote white pawn to queen",
			beforeFen: "8/3k3P/8/8/8/8/8/3K4 w - - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_H7).SetTargetSquare(types.SQUARE_H8).SetMoveType(move.PROMOTION).SetPromitionPieceType(types.QUEEN),
			afterFen:  "7Q/3k4/8/8/8/8/8/3K4 b - - 0 1",
		},
		{
			name:      "promote black pawn to rook",
			beforeFen: "8/3k4/8/8/8/8/6p1/3K4 b - - 0 1",
			m:         *new(move.Move).SetSourceSquare(types.SQUARE_G2).SetTargetSquare(types.SQUARE_G1).SetMoveType(move.PROMOTION).SetPromitionPieceType(types.ROOK),
			afterFen:  "8/3k4/8/8/8/8/8/3K2r1 w - - 0 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{}
			pos, err := NewFromFen(tt.beforeFen)
			assert.NoError(t, err)
			pos.MakeMove(tt.m, s)
			assert.Equal(t, tt.afterFen, pos.ToFen())
			pos.UnMakeMove(s)
			assert.Equal(t, tt.beforeFen, pos.ToFen())
		})
	}
}

func TestPosition_GeneratePseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name      string
		beforeFen string
		moves     string
	}{{
		name:      "Promotion",
		beforeFen: "r3k2r/p1ppqpb1/bn2pnp1/1N1PN3/1p2P3/5Q2/PPPB1PpP/R3KB1R b KQkq - 0 1",
		moves:     "a8b8 a8c8 a8d8 h8h2 h8h3 h8h4 h8h5 h8h6 h8h7 h8f8 h8g8 a6b5 a6b7 a6c8 g7h6 g7f8 e7c5 e7d6 e7d8 e7f8 b6a4 b6c4 b6d5 b6c8 f6e4 f6g4 f6d5 f6h5 f6h7 f6g8 g2g1n g2g1b g2g1r g2g1q g2f1n g2f1b g2f1r g2f1q g2h1n g2h1b g2h1r g2h1q b4b3 e6d5 g6g5 c7c5 c7c6 d7d6 e8g8 e8c8 e8d8 e8f8",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.beforeFen)
			assert.NoError(t, err)
			moves := move.NewMoveList()
			pos.GeneratePseudoLegalMoves(moves)
			assert.Equal(t, tt.moves, fmt.Sprint(moves))
		})
	}
}

func TestPosition_GeneratePseudoLegalCaptures(t *testing.T) {
	tests := []struct {
		name      string
		beforeFen string
		moves     string
	}{{
		name:      "Promotion",
		beforeFen: "r3k2r/p1ppqpb1/bn2pnp1/1N1PN3/1p2P3/5Q2/PPPB1PpP/R3KB1R b KQkq - 0 1",
		moves:     "h8h2 a6b5 b6d5 f6e4 f6d5 g2f1n g2f1b g2f1r g2f1q g2h1n g2h1b g2h1r g2h1q e6d5",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.beforeFen)
			assert.NoError(t, err)
			moves := move.NewMoveList()
			pos.GeneratePseudoLegalCaptures(moves)
			assert.Equal(t, tt.moves, fmt.Sprint(moves))
		})
	}
}
