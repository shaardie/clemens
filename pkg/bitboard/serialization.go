package bitboard

import "github.com/shaardie/clemens/pkg/types"

func AllSubnetsOf(b Bitboard) []Bitboard {
	setOfSubsets := []Bitboard{}
	subset := Empty
	for {
		setOfSubsets = append(setOfSubsets, subset)
		subset = (subset - b) & b
		if subset == Empty {
			break
		}
	}
	return setOfSubsets
}

func IsolatingSubsets(b Bitboard) []Bitboard {
	bs := []Bitboard{}
	for b != Empty {
		bs = append(bs, b&-b)
		b &= b - 1
	}
	return bs
}

func SquareIndexSerialization(b Bitboard) []int {
	idxs := [64]int{}
	length := 0
	for b != Empty {
		idxs[length] = LeastSignificantOneBit(b)
		length++
		b &= b - 1
	}
	return idxs[0:length]
}

func SquareIndexSerializationIterator(b Bitboard) func() int {
	return func() int {
		if b == Empty {
			return types.SQUARE_NONE
		}
		r := LeastSignificantOneBit(b)
		b &= b - 1
		return r
	}
}
