package bitboard

func IsolatingSubsets(b Bitboard) []Bitboard {
	bs := []Bitboard{}
	for b != Empty {
		bs = append(bs, b&-b)
		b &= b - 1
	}
	return bs
}
func SquareIndexSerialization(b Bitboard) []int {
	idxs := []int{}
	for b != Empty {
		idxs = append(idxs, LeastSignificantOneBit(b))
		b &= b - 1
	}
	return idxs
}
