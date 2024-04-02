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

func SquareIndexSerialization(b Bitboard) []uint8 {
	idxs := [64]uint8{}
	length := 0
	for b != Empty {
		idxs[length] = LeastSignificantOneBit(b)
		length++
		b &= b - 1
	}
	return idxs[0:length]
}

func SquareIndexSerializationNextSquare(b *Bitboard) uint8 {
	r := LeastSignificantOneBit(*b)
	*b &= *b - 1
	return r
}
