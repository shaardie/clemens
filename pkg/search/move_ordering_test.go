package search

import (
	"context"
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_orderMoves(t *testing.T) {
	pos, err := position.NewFromFen("r3k2r/p1ppqpb1/Bn4p1/3pN3/4nP2/2B5/PPP3QP/R3K2R w KQkq - 0 5")
	require.NoError(t, err)

	// Generate all moves and order them
	s := NewSearch(*pos)
	s.Search(context.TODO(), SearchParameter{Depth: 5})

	moves1 := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves1)
	s.orderMoves(pos, moves1, s.PV.GetBestMove(), move.NullMove, 0)

	moves2 := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves2)
	s.scoreMoves(pos, moves2, s.PV.GetBestMove(), move.NullMove, 0)
	for i := range moves2.Length() {
		moves2.SortIndex(i)
	}

	for i := range moves2.Length() {
		assert.Equal(t, moves1.Get(i).GetScore(), moves2.Get(i).GetScore())
	}
}
