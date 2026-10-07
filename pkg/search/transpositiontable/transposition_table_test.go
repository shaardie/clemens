package transpositiontable

import (
	"testing"

	"github.com/shaardie/clemens/pkg/evaluation"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGet_MateScorePlyAdjustment verifies that a node-relative mate score is
// adjusted correctly when read at a different ply.
// The stored value is expected to be node-relative ("mate in N plies from this
// node"); Get must convert it to a root-relative value for the current ply.
func TestGet_MateScorePlyAdjustment(t *testing.T) {
	hash := uint64(999)
	key := hash % numberOfBuckets
	tt[key] = bucket{}

	// Node-relative mate score: mate in 4 plies from this node.
	tt[key][0] = ttEntry{
		zobristHash:    hash,
		bestMove:       move.NullMove,
		score:          evaluation.INF - 4,
		depth:          10,
		ageAndNodeType: uint8(PVNode),
	}

	// Reading at ply 6 => mate at root-ply 10.
	score, use, _ := Get(hash, -evaluation.INF, evaluation.INF, 10, 6)
	require.True(t, use)
	assert.Equal(t, int16(evaluation.INF-10), score)

	// Reading at ply 0 => mate at root-ply 4.
	score, use, _ = Get(hash, -evaluation.INF, evaluation.INF, 10, 0)
	require.True(t, use)
	assert.Equal(t, int16(evaluation.INF-4), score)

	// Negative mate scores must be adjusted in the opposite direction.
	tt[key][0].score = -evaluation.INF + 4
	score, use, _ = Get(hash, -evaluation.INF, evaluation.INF, 10, 6)
	require.True(t, use)
	assert.Equal(t, int16(-evaluation.INF+10), score)
}

// TestTT_MateScoreStoreGetRoundTrip verifies the full store + get round-trip:
// the search passes a root-relative mate score together with the ply, and the
// stored node-relative value must convert back correctly for any ply.
func TestTT_MateScoreStoreGetRoundTrip(t *testing.T) {
	hash := uint64(1234)
	key := hash % numberOfBuckets
	tt[key] = bucket{}

	// Search found a mate at root-ply 8 while at node ply 5 (mate in 3).
	PotentiallySave(hash, move.NullMove, 10, 5, evaluation.INF-8, PVNode, 0)

	// Reading at the same ply returns the original root-relative value.
	score, use, _ := Get(hash, -evaluation.INF, evaluation.INF, 10, 5)
	require.True(t, use)
	assert.Equal(t, int16(evaluation.INF-8), score)

	// Reading at a different ply (transposition): mate distance 3 => mate at ply 6.
	score, use, _ = Get(hash, -evaluation.INF, evaluation.INF, 10, 3)
	require.True(t, use)
	assert.Equal(t, int16(evaluation.INF-6), score)
}
