package nnue

import (
	"github.com/shaardie/clemens/pkg/types"
)

type Accumulator [types.COLOR_NUMBER][HiddenSize]float32

// NewAccumulator creates a new accumulator.
// Initialize the bias layer
func NewAccumulator() Accumulator {
	return Accumulator{
		types.WHITE: m.FTBias,
		types.BLACK: m.FTBias,
	}
}

// Refresh calculates the accumulator completly from the existing features.
// whiteFeatures/blackFeatures: slices of active features indices (0–767).
func (acc *Accumulator) Refresh(whiteFeatures, blackFeatures []int) {
	acc[types.WHITE] = m.FTBias
	acc[types.BLACK] = m.FTBias

	for _, f := range whiteFeatures {
		for j := range HiddenSize {
			acc[types.WHITE][j] += m.FTWeight[f][j]
		}
	}
	for _, f := range blackFeatures {
		for j := range HiddenSize {
			acc[types.BLACK][j] += m.FTWeight[f][j]
		}
	}
}

// Activate a feature in the accumulator (add a piece on a square)
// whiteFeature/blackFeature: Feature index from whites/blacks perspective
func (acc *Accumulator) Activate(whiteFeature, blackFeature int) {
	for j := range HiddenSize {
		acc[types.WHITE][j] += m.FTWeight[whiteFeature][j]
		acc[types.BLACK][j] += m.FTWeight[blackFeature][j]
	}
}

// Deactivate a feature in the accumulator (remove a piece from a square)
// whiteFeature/blackFeature: Feature index from whites/blacks perspective
func (acc *Accumulator) Deactivate(whiteFeature, blackFeature int) {
	for j := range HiddenSize {
		acc[types.WHITE][j] -= m.FTWeight[whiteFeature][j]
		acc[types.BLACK][j] -= m.FTWeight[blackFeature][j]
	}
}

func crelu(x float32) float32 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	return x
}

// Evaluate the value in centipawns
// aus Sicht der Seite am Zug (positiv = gut für STM).
func (acc *Accumulator) Evaluate(c types.Color) int16 {
	// Perspektive bestimmen: STM zuerst, dann Gegner
	var stm, opp *[HiddenSize]float32
	if c == types.WHITE {
		stm = &acc[types.WHITE]
		opp = &acc[types.BLACK]
	} else {
		stm = &acc[types.BLACK]
		opp = &acc[types.WHITE]
	}

	// Clipped ReLU on accumulator values
	var input [HiddenSize * 2]float32
	for i := range HiddenSize {
		input[i] = crelu(stm[i])
	}
	for i := range HiddenSize {
		input[HiddenSize+i] = crelu(opp[i])
	}

	// layer 1: 512 → 32
	var l1 [L1Size]float32
	for i := range L1Size {
		sum := m.L1Bias[i]
		for j := 0; j < HiddenSize*2; j++ {
			sum += input[j] * m.L1Weight[i][j]
		}
		l1[i] = crelu(sum)
	}

	// layer 2: 32 → 32
	var l2 [L2Size]float32
	for i := range L2Size {
		sum := m.L2Bias[i]
		for j := range L1Size {
			sum += l1[j] * m.L2Weight[i][j]
		}
		l2[i] = crelu(sum)
	}

	// layer 3: 32 → 1
	out := m.OutBias
	for i := range L2Size {
		out += l2[i] * m.OutWeight[i]
	}

	// calculate in centipawns:
	return int16(EvalScale * float64(out))
}
