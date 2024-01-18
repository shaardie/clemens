package bitboard

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
	idxs := make([]int, 0, b.PopulationCount())
	for b != Empty {
		idxs = append(idxs, LeastSignificantOneBit(b))
		b &= b - 1
	}
	return idxs
}
