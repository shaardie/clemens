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
	assert.Equal(t, types.SQUARE_H7, m.GetSourceSquare())
	assert.Equal(t, types.SQUARE_H8, m.GetTargetSquare())
	assert.Equal(t, PROMOTION, m.GetMoveType())
	assert.Equal(t, types.ROOK, m.GetPromitionPieceType())
}
