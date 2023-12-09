package bitboard

func Equal(b1, b2 Bitboard) bool {
	return b1 == b2
}

func Intersection(b1, b2 Bitboard) Bitboard {
	return b1 & b2
}

func Union(b1, b2 Bitboard) Bitboard {
	return b1 | b2
}

func (b Bitboard) Complement() Bitboard {
	return ^b
}

func Difference(b1, b2 Bitboard) Bitboard {
	return b1 & ^b2
}

func Implication(b1, b2 Bitboard) Bitboard {
	return ^b1 | b2
}

func SymmetricDifference(b1, b2 Bitboard) Bitboard {
	return b1 ^ b2
}

func Equivalence(b1, b2 Bitboard) Bitboard {
	return ^(b1 ^ b2)
}

func Majority(b1, b2, b3 Bitboard) Bitboard {
	return (b1 & b2) | (b2 & b3) | (b1 & b3)
}

func GreaterOne(bs ...Bitboard) (r Bitboard) {
	for i := 1; i < len(bs); i++ {
		u := Empty
		for j := 0; j < i; j++ {
			u |= bs[j]
		}
		r |= bs[i] & u
	}
	return
}
