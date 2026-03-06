package position

import (
	"github.com/shaardie/clemens/pkg/nnue"
	"github.com/shaardie/clemens/pkg/types"
)

func (pos *Position) initAccumulator() {
	features := [types.COLOR_NUMBER][]int{
		make([]int, 0, types.SQUARE_NUMBER),
		make([]int, 0, types.SQUARE_NUMBER),
	}
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		piece := pos.GetPiece(square)
		if piece == types.NO_PIECE {
			continue
		}
		whiteIndex, blackIndex := nnue.FeatureIndex(piece, square)
		features[types.WHITE] = append(features[types.WHITE], whiteIndex)
		features[types.BLACK] = append(features[types.BLACK], blackIndex)
	}
	pos.Accumulator = nnue.NewAccumulator()
	pos.Accumulator.Refresh(features[types.WHITE], features[types.BLACK])
}
