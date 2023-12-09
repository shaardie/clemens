package bitboard

func SouthOne(b Bitboard) Bitboard     { return b >> 8 }
func NorthOne(b Bitboard) Bitboard     { return b << 8 }
func EastOne(b Bitboard) Bitboard      { return (b << 1) & notAFile }
func NorthEastOne(b Bitboard) Bitboard { return (b << 9) & notAFile }
func SouthEastOne(b Bitboard) Bitboard { return (b >> 7) & notAFile }
func WestOne(b Bitboard) Bitboard      { return (b >> 1) & notHFile }
func SouthWestOne(b Bitboard) Bitboard { return (b >> 9) & notHFile }
func NorthWest(b Bitboard) Bitboard    { return (b << 7) & notHFile }

func RotateLeft(b Bitboard, s int) Bitboard  { return (b << s) | (b >> (64 - s)) }
func RotateRight(b Bitboard, s int) Bitboard { return (b >> s) | (b << (64 - s)) }

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
