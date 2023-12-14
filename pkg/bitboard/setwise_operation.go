package bitboard

import "math/bits"

func Equal(b1, b2 Bitboard) bool {
	return b1 == b2
}

func Intersection(b1, b2 Bitboard) Bitboard {
	return b1 & b2
}

func Union(b1, b2 Bitboard) Bitboard {
	return b1 | b2
}

func Complement(b Bitboard) Bitboard {
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

func SouthOne(b Bitboard) Bitboard     { return b >> 8 }
func NorthOne(b Bitboard) Bitboard     { return b << 8 }
func EastOne(b Bitboard) Bitboard      { return (b << 1) & notAFile }
func NorthEastOne(b Bitboard) Bitboard { return (b << 9) & notAFile }
func SouthEastOne(b Bitboard) Bitboard { return (b >> 7) & notAFile }
func WestOne(b Bitboard) Bitboard      { return (b >> 1) & notHFile }
func SouthWestOne(b Bitboard) Bitboard { return (b >> 9) & notHFile }
func NorthWestOne(b Bitboard) Bitboard { return (b << 7) & notHFile }

func RotateLeft(b Bitboard, s int) Bitboard {
	return Bitboard(bits.RotateLeft64(uint64(b), s))
}
func RotateRight(b Bitboard, s int) Bitboard {
	return Bitboard(bits.RotateLeft64(uint64(b), -s))
}

func GeneralShift(b Bitboard, s int) Bitboard {
	switch {
	case s < 0:
		return (b >> -s)
	case s > 0:
		return (b << s)
	default:
		return b
	}
}

func BitBySquare(square uint64) Bitboard { return 1 << square }
