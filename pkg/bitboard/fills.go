package bitboard

func NorthFill(b Bitboard) Bitboard {
	b |= (b << 8)
	b |= (b << 16)
	b |= (b << 32)
	return b
}

func SouthFill(b Bitboard) Bitboard {
	b |= (b >> 8)
	b |= (b >> 16)
	b |= (b >> 32)
	return b
}

func FileFill(b Bitboard) Bitboard {
	return NorthFill(b) | SouthFill(b)
}
