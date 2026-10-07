package position

import (
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPosition_ZobristHash(t *testing.T) {
	pos := New()
	assert.NotZero(t, pos.ZobristHash)

	var m move.Move
	m.SetSourceSquare(types.SQUARE_E2)
	m.SetTargetSquare(types.SQUARE_E4)
	pos.MakeMove(m)
	afterMoveHash := pos.ZobristHash
	pos.initZobristHash()
	reinitHash := pos.ZobristHash
	assert.Equal(t, reinitHash, afterMoveHash)
}

// TestPosition_ZobristHashConsistency verifies that the incrementally updated
// zobrist hash matches the hash of the same position re-initialized from FEN.
// This is the key property for the transposition table and repetition detection
// and catches bugs in the incremental castling/en-passant updates.
func TestPosition_ZobristHashConsistency(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		moves []string
	}{
		{
			name:  "startpos",
			fen:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []string{"e2e4", "e7e5", "g1f3", "b8c6", "f1c4"},
		},
		{
			name:  "both castle",
			fen:   "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			moves: []string{"e1g1", "e8g8"},
		},
		{
			name:  "rooks move away",
			fen:   "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			moves: []string{"a1a2", "a8a7", "h1h2", "h8h7"},
		},
		{
			name:  "only one right then king moves",
			fen:   "r3k2r/8/8/8/8/8/8/R3K2R w Kkq - 0 1",
			moves: []string{"e1e2", "e8e7"},
		},
		{
			name:  "king and rook return",
			fen:   "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			moves: []string{"e1e2", "e2e1", "a1a2", "a2a1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.fen)
			require.NoError(t, err)

			for _, mv := range tt.moves {
				err := pos.MakeMoveFromString(mv)
				require.NoError(t, err)

				reinit, err := NewFromFen(pos.ToFen())
				require.NoError(t, err)
				assert.Equal(t, reinit.ZobristHash, pos.ZobristHash,
					"zobrist hash after %q differs from re-initialized hash", mv)
			}
		})
	}
}
