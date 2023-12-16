package bitboard

func RankOfSquare(square int) int {
	return square >> 3
}

func RankMaskOfSquare(square int) Bitboard {
	return RankMask1 << (8 * RankOfSquare(square))
}

func FileOfSquare(square int) int {
	return square & 7
}

func SquareFromRankAndFile(rank int, file int) int {
	return (rank << 3) + file
}

func BitBySquares(squares ...int) Bitboard {
	b := Empty
	for _, s := range squares {
		b |= One << s
	}
	return b
}

func BitboardFromRankAndFile(rank int, file int) Bitboard {
	return BitBySquares(SquareFromRankAndFile(rank, file))
}
