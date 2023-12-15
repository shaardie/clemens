package main

import (
	"fmt"

	"github.com/shaardie/clemens/pkg/bitboard"
)

func main() {
	fmt.Println(bitboard.Bitboard(0x00000000000000ff).PrettyString())
}
