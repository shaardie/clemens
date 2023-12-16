package main

import (
	"fmt"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/pieces"
)

func main() {
	fmt.Println(pieces.SlidingAttacks(pieces.RookType, bitboard.SQUARE_F3, 0).PrettyString())
}
