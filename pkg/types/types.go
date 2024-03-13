package types

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	SQUARE_A1 int = iota
	SQUARE_B1
	SQUARE_C1
	SQUARE_D1
	SQUARE_E1
	SQUARE_F1
	SQUARE_G1
	SQUARE_H1
	SQUARE_A2
	SQUARE_B2
	SQUARE_C2
	SQUARE_D2
	SQUARE_E2
	SQUARE_F2
	SQUARE_G2
	SQUARE_H2
	SQUARE_A3
	SQUARE_B3
	SQUARE_C3
	SQUARE_D3
	SQUARE_E3
	SQUARE_F3
	SQUARE_G3
	SQUARE_H3
	SQUARE_A4
	SQUARE_B4
	SQUARE_C4
	SQUARE_D4
	SQUARE_E4
	SQUARE_F4
	SQUARE_G4
	SQUARE_H4
	SQUARE_A5
	SQUARE_B5
	SQUARE_C5
	SQUARE_D5
	SQUARE_E5
	SQUARE_F5
	SQUARE_G5
	SQUARE_H5
	SQUARE_A6
	SQUARE_B6
	SQUARE_C6
	SQUARE_D6
	SQUARE_E6
	SQUARE_F6
	SQUARE_G6
	SQUARE_H6
	SQUARE_A7
	SQUARE_B7
	SQUARE_C7
	SQUARE_D7
	SQUARE_E7
	SQUARE_F7
	SQUARE_G7
	SQUARE_H7
	SQUARE_A8
	SQUARE_B8
	SQUARE_C8
	SQUARE_D8
	SQUARE_E8
	SQUARE_F8
	SQUARE_G8
	SQUARE_H8
	SQUARE_NUMBER
	SQUARE_NONE = 64
)

const (
	fileToChar string = "abcdefgh"
)

func RankOfSquare(square int) int {
	return square >> 3
}

func FileOfSquare(square int) int {
	return square & 7
}

func SquareFromRankAndFile(rank int, file int) int {
	return (rank << 3) + file
}

func SquareToString(square int) string {
	return fmt.Sprintf(
		"%s%d",
		string(fileToChar[FileOfSquare(square)]),
		RankOfSquare(square)+1)
}

func SquareFromString(square string) (int, error) {
	var file, rank int
	for i, r := range square {
		switch i {
		case 0:
			file = strings.IndexRune(fileToChar, r)
			if file == -1 {
				return 0, fmt.Errorf("failed to get file")
			}
		case 1:
			if !unicode.IsDigit(r) {
				return 0, fmt.Errorf("failed to get rank")
			}
			rank = int(r-'0') - 1
		default:
			return 0, fmt.Errorf("token to long")
		}
	}
	return SquareFromRankAndFile(rank, file), nil
}

const (
	RANK_1 int = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
	RANK_NUMBER
)

const (
	FILE_A int = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
	FILE_NUMBER
)

type PieceType int

const (
	PAWN PieceType = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	PIECE_TYPE_NUMBER
)

func (p PieceType) String() string {
	switch p {
	case KNIGHT:
		return "n"
	case BISHOP:
		return "b"
	case ROOK:
		return "r"
	case QUEEN:
		return "q"
	}
	return ""
}

func PieceTypeFromString(s string) (PieceType, error) {
	switch s {
	case "p":
		return PAWN, nil
	case "n":
		return KNIGHT, nil
	case "b":
		return BISHOP, nil
	case "r":
		return ROOK, nil
	case "q":
		return QUEEN, nil
	default:
		return 0, errors.New("unknown piece type")
	}
}

type Color int

const (
	WHITE Color = iota
	BLACK
	COLOR_NUMBER
)

func SwitchColor(c Color) Color {
	if c == BLACK {
		return WHITE
	}
	return BLACK
}

type Piece int

const (
	NO_PIECE Piece = iota
	WHITE_PAWN
	WHITE_KNIGHT
	WHITE_BISHOP
	WHITE_ROOK
	WHITE_QUEEN
	WHITE_KING
)
const (
	BLACK_PAWN Piece = iota + WHITE_PAWN + 8
	BLACK_KNIGHT
	BLACK_BISHOP
	BLACK_ROOK
	BLACK_QUEEN
	BLACK_KING
)

func (p Piece) Color() Color {
	return Color(p >> 3)
}

func (p Piece) Type() PieceType {
	return PieceType((p & 7) - 1)
}

const pieceToChar string = " PNBRQK  pnbrqk"

func (p Piece) ToChar() rune {
	for idx, r := range pieceToChar {
		if idx == int(p) {
			return r
		}
	}
	panic("rune not found")
}

func NewPiece(c Color, pt PieceType) Piece {
	return Piece((int(pt) + 1) + int(c)*8)
}

func NewPieceFromChar(r rune) (Piece, error) {
	idx := strings.IndexRune(pieceToChar, r)
	if idx == -1 {
		return 0, fmt.Errorf("%v is no valid piece", r)
	}
	return Piece(idx), nil
}
