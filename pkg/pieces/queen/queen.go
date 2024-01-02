package queen

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/pieces/utils"
)

// AttacksBySquare returns the attacks for a given square.
func AttacksBySquare(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	return rook.AttacksBySquare(square, occupied) | bishop.AttacksBySquare(square, occupied)
}

// AttacksByBitboard calculates the AttacksByBitboard of the queen for the given square and occupation
func AttacksByBitboard(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	return bishop.AttacksByBitboard(square, occupied) | rook.AttacksByBitboard(square, occupied)
}

// GenerateMoves generates all moves for all squares to all destinations by a given occupation
func GenerateMoves(squares, occupied, destinations bitboard.Bitboard) []move.Move {
	return utils.GenerateMoves(
		squares,
		occupied,
		destinations,
		AttacksByBitboard,
	)
}
