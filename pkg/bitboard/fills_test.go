package bitboard

import (
	"testing"

	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNorthFill(t *testing.T) {
	assert.Equal(
		t,
		NorthFill(BitBySquares(types.SQUARE_A4)),
		BitBySquares(
			types.SQUARE_A4,
			types.SQUARE_A5,
			types.SQUARE_A6,
			types.SQUARE_A7,
			types.SQUARE_A8,
		),
	)
}

func TestSouthFill(t *testing.T) {
	assert.Equal(
		t,
		SouthFill(BitBySquares(types.SQUARE_A4)),
		BitBySquares(
			types.SQUARE_A4,
			types.SQUARE_A3,
			types.SQUARE_A2,
			types.SQUARE_A1,
		),
	)
}

func TestFileFill(t *testing.T) {
	assert.Equal(
		t,
		FileFill(BitBySquares(types.SQUARE_C3)),
		FileMaskC,
	)
}
