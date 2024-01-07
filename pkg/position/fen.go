package position

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	fileToChar string = "abcdefgh"
)

// NewFromFen creates a new position from a FEN string, see https://www.chessprogramming.org/Forsyth-Edwards_Notation#En_passant_target_square
func NewFromFen(fen string) (*Position, error) {
	pos := &Position{}

	tokens := strings.Split(fen, " ")
	if len(tokens) != 6 {
		return nil, fmt.Errorf("wrong number of tokens %v", len(tokens))
	}
	err := pos.fenSetPieces(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("failed to set pieces from fen token, %w", err)
	}

	err = pos.fenSetSideToMove(tokens[1])
	if err != nil {
		return nil, fmt.Errorf("failed to set side to move from fen token, %w", err)
	}

	// Set castling
	err = pos.fenSetCastling(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("failed to set castling from fen token, %w", err)
	}

	// Set castling
	err = pos.fenSetEnPassant(tokens[3])
	if err != nil {
		return nil, fmt.Errorf("failed to set en passant from fen token, %w", err)
	}

	pos.halfMoveClock, err = strconv.Atoi(tokens[4])
	if err != nil {
		return nil, fmt.Errorf("failed to set half move clock from fen token, %w", err)
	}

	pos.numberOfFullMoves, err = strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("failed to set number of full moves fen token, %w", err)
	}

	return pos, nil
}

// ToFen creates FEN string from position, see https://www.chessprogramming.org/Forsyth-Edwards_Notation
func (pos *Position) ToFen() string {
	sb := strings.Builder{}
	for rank := types.RANK_8; rank >= types.RANK_1; rank = rank - 1 {
		for file := types.FILE_A; file <= types.FILE_H; file++ {
			emptyCount := 0
			for file <= types.FILE_H && pos.Empty(bitboard.SquareFromRankAndFile(rank, file)) {
				emptyCount++
				file++
			}
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
			}

			if file <= types.FILE_H {
				sb.WriteRune(
					pos.GetPiece(bitboard.SquareFromRankAndFile(rank, file)).ToChar(),
				)
			}
		}
		if rank > types.RANK_1 {
			sb.WriteRune('/')
		}
	}
	if pos.sideToMove == types.WHITE {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	if !pos.CanCastle(ANY_CASTLING) {
		sb.WriteRune('-')
	} else {
		if pos.CanCastle(WHITE_CASTLING_KING) {
			sb.WriteRune('K')
		}
		if pos.CanCastle(WHITE_CASTLING_QUEEN) {
			sb.WriteRune('Q')
		}
		if pos.CanCastle(BLACK_CASTLING_KING) {
			sb.WriteRune('k')
		}
		if pos.CanCastle(BLACK_CASTLING_QUEEN) {
			sb.WriteRune('q')
		}
	}
	sb.WriteRune(' ')

	if pos.enPassant == types.SQUARE_NONE {
		sb.WriteRune('-')
	} else {
		sb.WriteByte(fileToChar[bitboard.FileOfSquare(pos.enPassant)])
		sb.WriteString(strconv.Itoa(bitboard.RankOfSquare(pos.enPassant)))
	}
	sb.WriteRune(' ')

	sb.WriteString(strconv.Itoa(bitboard.RankOfSquare(pos.halfMoveClock)))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(pos.numberOfFullMoves))
	return sb.String()
}

// fenSetPieces set piece positions from part of the fen string
func (pos *Position) fenSetPieces(token string) error {
	reader := strings.NewReader(token)
	square := types.SQUARE_A8
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("failed to read piece placement, %w", err)
		}
		if unicode.IsDigit(r) {
			// Jump forward in file
			square += int(r - '0')
			continue
		} else if r == '/' {
			// Jump to the beginning of the previous rank
			square -= 2 * 8
			continue
		}

		// Set piece
		p, err := types.PieceFromChar(r)
		if err != nil {
			return fmt.Errorf("failed to get piece from rune, %w", err)
		}
		pos.SetPiece(p, square)
		square++
	}
}

// fenSetSideToMove set piece positions from part of the fen string
func (pos *Position) fenSetSideToMove(token string) error {
	if len(token) != 1 {
		return fmt.Errorf("fen token for side to move has wrong number of chars")
	}

	if token == "w" {
		pos.sideToMove = types.WHITE
		return nil
	}
	if token == "b" {
		pos.sideToMove = types.BLACK
		return nil
	}

	return fmt.Errorf("unknown color for side to move")
}

// fenSetCastling set castling from part of the fen string
func (pos *Position) fenSetCastling(token string) error {
	pos.castling = NO_CASTLING

	// No castling
	if token == "-" {
		return nil
	}

	reader := strings.NewReader(token)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("failed to read castling, %w", err)
		}
		if r == 'K' {
			pos.castling |= WHITE_CASTLING_KING
		} else if r == 'Q' {
			pos.castling |= WHITE_CASTLING_QUEEN
		} else if r == 'k' {
			pos.castling |= BLACK_CASTLING_KING
		} else if r == 'q' {
			pos.castling |= BLACK_CASTLING_QUEEN
		} else {
			return fmt.Errorf("unknown castling rule, %v", string(r))
		}
	}
}

// fenSetEnPassant set en passant from part of the fen string
func (pos *Position) fenSetEnPassant(token string) error {
	pos.enPassant = types.SQUARE_NONE

	// No en passant
	if token == "-" {
		return nil
	}

	var file, rank int
	for i, r := range token {
		switch i {
		case 0:
			file = strings.IndexRune(fileToChar, r)
			if file == -1 {
				return fmt.Errorf("failed to get file")
			}
		case 1:
			if !unicode.IsDigit(r) {
				return fmt.Errorf("failed to get rank")
			}
			rank = int(r-'0') - 1
		default:
			return fmt.Errorf("token to long")
		}
	}

	pos.enPassant = bitboard.SquareFromRankAndFile(rank, file)
	return nil
}
