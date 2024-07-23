package nnue

import (
	"github.com/shaardie/clemens/pkg/types"
)

type Accumulator [types.COLOR_NUMBER][_M]float64

func (n *Accumulator) Refresh(activeFeatures []int, color types.Color) {
	// First we copy the layer bias, that's our starting point
	copy(m.L0.Bias, n[color][:])

	// Then we just accumulate all the columns for the active features. That's what accumulators do!
	for _, feature := range activeFeatures {
		for i := range n[color] {
			n[color][i] += m.L0.Weights[i][feature]
		}
	}
}

func (n *Accumulator) Update(addedFeatures, removedFeatures []int, color types.Color) {
	// Then we subtract the weights of the removed features
	for _, feature := range removedFeatures {
		for i := range n[color] {
			n[color][i] -= m.L0.Weights[i][feature]
		}
	}

	// Similar for the added features, but add instead of subtracting
	for _, feature := range addedFeatures {
		for i := range n[color] {
			n[color][i] += m.L0.Weights[i][feature]
		}
	}
}

func (n *Accumulator) Evaluate(we types.Color) int16 {
	var input [2 * _M]float64
	them := types.SwitchColor(we)
	for i := range _M {
		input[i] = n[we][i]
		input[_M+i] = n[them][i]
	}

	var creluL0Output [2 * _M]float64
	crelu(input[:], creluL0Output[:])

	var l1LinearOutput [_N]float64
	linear(&m.L1, creluL0Output[:], l1LinearOutput[:])

	var creluL1Output [_N]float64
	crelu(l1LinearOutput[:], creluL1Output[:])

	var l2LinearOutput [_K]float64
	linear(&m.L1, creluL1Output[:], l2LinearOutput[:])

	return int16(10000 * l2LinearOutput[0])
}
