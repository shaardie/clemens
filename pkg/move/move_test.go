package move

import (
	"testing"

	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	var m Move
	m.SetSourceSquare(types.SQUARE_H7)
	m.SetTargetSquare(types.SQUARE_H8)
	m.SetMoveType(PROMOTION)
	m.SetPromitionPieceType(types.ROOK)
	m.SetScore(1234)
	assert.Equal(t, types.SQUARE_H7, m.GetSourceSquare())
	assert.Equal(t, types.SQUARE_H8, m.GetTargetSquare())
	assert.Equal(t, PROMOTION, m.GetMoveType())
	assert.Equal(t, types.ROOK, m.GetPromitionPieceType())
	assert.Equal(t, uint16(1234), m.GetScore())
}

// Ensure that Moves are considered equal, even if they have different scores
func TestMoveCompare(t *testing.T) {
	var m1, m2 Move
	m1.SetSourceSquare(types.SQUARE_H7)
	m1.SetTargetSquare(types.SQUARE_H8)
	m1.SetMoveType(PROMOTION)
	m1.SetPromitionPieceType(types.ROOK)
	m2 = m1
	m1.SetScore(1234)
	m2.SetScore(4321)
	assert.NotEqual(t, m2, m1, "m1 and m2 should not be equal in bits, because they have different scores")
	assert.True(t, m1.Equal(m2), "m1 and m2 should be equal")
}

// Ensure that scores are properly idempotent
func TestSetScore(t *testing.T) {
	var m Move

	for _, score := range []uint16{0x00FF, 0xFF00} {
		m.SetScore(score)
		assert.Equal(t, score, m.GetScore(), "wrong move score %x", score)
	}
}
