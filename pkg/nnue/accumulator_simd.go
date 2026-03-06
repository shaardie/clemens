//go:build goexperiment.simd

package nnue

import (
	"simd/archsimd"

	"github.com/shaardie/clemens/pkg/types"
)

// Compile-time check: HiddenSize must be divisible by 8 for SIMD.
var _ [0]struct {
	_ [HiddenSize%8 | (HiddenSize*2)%8]byte
}

type Accumulator [types.COLOR_NUMBER][HiddenSize]float32

// simdOnes is a precomputed vector of 1.0s for Clipped ReLU upper bound.
var simdOnes = archsimd.LoadFloat32x8Slice(
	[]float32{1, 1, 1, 1, 1, 1, 1, 1},
)

// simdZeros is the zero vector for Clipped ReLU lower bound.
var simdZeros = archsimd.Float32x8{}

func NewAccumulator() Accumulator {
	return Accumulator{
		types.WHITE: m.FTBias,
		types.BLACK: m.FTBias,
	}
}

func (acc *Accumulator) Refresh(whiteFeatures, blackFeatures []int) {
	for j := 0; j < HiddenSize; j += 8 {
		b := archsimd.LoadFloat32x8Slice(m.FTBias[j:])
		b.StoreSlice(acc[types.WHITE][j:])
		b.StoreSlice(acc[types.BLACK][j:])
	}

	for _, f := range whiteFeatures {
		w := m.FTWeight[f][:]
		for j := 0; j < HiddenSize; j += 8 {
			a := archsimd.LoadFloat32x8Slice(acc[types.WHITE][j:])
			wv := archsimd.LoadFloat32x8Slice(w[j:])
			a.Add(wv).StoreSlice(acc[types.WHITE][j:])
		}
	}
	for _, f := range blackFeatures {
		w := m.FTWeight[f][:]
		for j := 0; j < HiddenSize; j += 8 {
			a := archsimd.LoadFloat32x8Slice(acc[types.BLACK][j:])
			wv := archsimd.LoadFloat32x8Slice(w[j:])
			a.Add(wv).StoreSlice(acc[types.BLACK][j:])
		}
	}
}

func (acc *Accumulator) Activate(whiteFeature, blackFeature int) {
	ww := m.FTWeight[whiteFeature][:]
	bw := m.FTWeight[blackFeature][:]
	for j := 0; j < HiddenSize; j += 8 {
		a := archsimd.LoadFloat32x8Slice(acc[types.WHITE][j:])
		w := archsimd.LoadFloat32x8Slice(ww[j:])
		a.Add(w).StoreSlice(acc[types.WHITE][j:])

		a = archsimd.LoadFloat32x8Slice(acc[types.BLACK][j:])
		w = archsimd.LoadFloat32x8Slice(bw[j:])
		a.Add(w).StoreSlice(acc[types.BLACK][j:])
	}
}

func (acc *Accumulator) Deactivate(whiteFeature, blackFeature int) {
	ww := m.FTWeight[whiteFeature][:]
	bw := m.FTWeight[blackFeature][:]
	for j := 0; j < HiddenSize; j += 8 {
		a := archsimd.LoadFloat32x8Slice(acc[types.WHITE][j:])
		w := archsimd.LoadFloat32x8Slice(ww[j:])
		a.Sub(w).StoreSlice(acc[types.WHITE][j:])

		a = archsimd.LoadFloat32x8Slice(acc[types.BLACK][j:])
		w = archsimd.LoadFloat32x8Slice(bw[j:])
		a.Sub(w).StoreSlice(acc[types.BLACK][j:])
	}
}

// hsum reduces a Float32x8 to a single float32 by summing all 8 elements.
func hsum(v archsimd.Float32x8) float32 {
	sum4 := v.GetHi().Add(v.GetLo())
	return sum4.GetElem(0) + sum4.GetElem(1) + sum4.GetElem(2) + sum4.GetElem(3)
}

func (acc *Accumulator) Evaluate(c types.Color) int16 {
	var stm, opp *[HiddenSize]float32
	if c == types.WHITE {
		stm = &acc[types.WHITE]
		opp = &acc[types.BLACK]
	} else {
		stm = &acc[types.BLACK]
		opp = &acc[types.WHITE]
	}

	// Apply Clipped ReLU using SIMD: clamp(x, 0, 1)
	var input [HiddenSize * 2]float32
	for i := 0; i < HiddenSize; i += 8 {
		v := archsimd.LoadFloat32x8Slice(stm[i:])
		v = v.Max(simdZeros).Min(simdOnes)
		v.StoreSlice(input[i:])
	}
	for i := 0; i < HiddenSize; i += 8 {
		v := archsimd.LoadFloat32x8Slice(opp[i:])
		v = v.Max(simdZeros).Min(simdOnes)
		v.StoreSlice(input[HiddenSize+i:])
	}

	// Layer 1: (HiddenSize*2) → L1Size
	inputSlice := input[:]
	var l1 [L1Size]float32
	for i := range L1Size {
		weights := m.L1Weight[i][:]
		acc0 := archsimd.Float32x8{}
		for j := 0; j < HiddenSize*2; j += 8 {
			inp := archsimd.LoadFloat32x8Slice(inputSlice[j:])
			w := archsimd.LoadFloat32x8Slice(weights[j:])
			acc0 = inp.MulAdd(w, acc0)
		}
		l1[i] = crelu(hsum(acc0) + m.L1Bias[i])
	}

	// Layer 2: L1Size → L2Size
	var l2 [L2Size]float32
	for i := range L2Size {
		sum := m.L2Bias[i]
		for j := range L1Size {
			sum += l1[j] * m.L2Weight[i][j]
		}
		l2[i] = crelu(sum)
	}

	// Output: L2Size → 1
	out := m.OutBias
	for i := range L2Size {
		out += l2[i] * m.OutWeight[i]
	}

	return int16(EvalScale * float64(out))
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
