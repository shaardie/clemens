package position

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/shaardie/clemens/pkg/types"
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

	halfMoveClock, err := strconv.Atoi(tokens[4])
	if err != nil {
		return nil, fmt.Errorf("failed to set half move clock from fen token, %w", err)
	}
	pos.HalfMoveClock = uint8(halfMoveClock)

	numberOfFullMoves, err := strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("failed to set number of full moves fen token, %w", err)
	}
	pos.Ply = uint8(2*numberOfFullMoves - 1)
	if pos.SideToMove == types.WHITE {
		pos.Ply--
	}

	// Create initial zobrist hash
	pos.initZobristHash()

	// Generate Helper Bitboards
	pos.generateHelperBitboards()

	return pos, nil
}

// ToFen creates FEN string from position, see https://www.chessprogramming.org/Forsyth-Edwards_Notation
func (pos *Position) ToFen() string {
	sb := strings.Builder{}
	rank := types.RANK_8
	for {
		for file := types.FILE_A; file <= types.FILE_H; file++ {
			emptyCount := 0
			for file <= types.FILE_H && pos.Empty(types.SquareFromRankAndFile(rank, file)) {
				emptyCount++
				file++
			}
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
			}

			if file <= types.FILE_H {
				sb.WriteRune(
					pos.GetPiece(types.SquareFromRankAndFile(rank, file)).ToChar(),
				)
			}
		}
		if rank > types.RANK_1 {
			sb.WriteRune('/')
		}
		if rank == types.RANK_1 {
			break
		}
		rank--
	}
	if pos.SideToMove == types.WHITE {
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

	if pos.EnPassant == types.SQUARE_NONE {
		sb.WriteRune('-')
	} else {
		sb.WriteString(types.SquareToString(pos.EnPassant))
	}
	sb.WriteRune(' ')

	sb.WriteString(strconv.Itoa(int(pos.HalfMoveClock)))
	sb.WriteRune(' ')

	numberOfFullMoves := int(pos.Ply/2 + 1)
	sb.WriteString(strconv.Itoa(numberOfFullMoves))
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
			square += uint8(r - '0')
			continue
		} else if r == '/' {
			// Jump to the beginning of the previous rank
			square -= 2 * 8
			continue
		}

		// Set piece
		p, err := types.NewPieceFromChar(r)
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
		pos.SideToMove = types.WHITE
		return nil
	}
	if token == "b" {
		pos.SideToMove = types.BLACK
		return nil
	}

	return fmt.Errorf("unknown color for side to move")
}

// fenSetCastling set castling from part of the fen string
func (pos *Position) fenSetCastling(token string) error {
	pos.Castling = NO_CASTLING

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
			pos.Castling |= WHITE_CASTLING_KING
		} else if r == 'Q' {
			pos.Castling |= WHITE_CASTLING_QUEEN
		} else if r == 'k' {
			pos.Castling |= BLACK_CASTLING_KING
		} else if r == 'q' {
			pos.Castling |= BLACK_CASTLING_QUEEN
		} else {
			return fmt.Errorf("unknown castling rule, %v", string(r))
		}
	}
}

// fenSetEnPassant set en passant from part of the fen string
func (pos *Position) fenSetEnPassant(token string) error {
	pos.EnPassant = types.SQUARE_NONE

	// No en passant
	if token == "-" {
		return nil
	}

	square, err := types.SquareFromString(token)
	if err != nil {
		return fmt.Errorf("failed to parse token as square, %w", err)
	}

	pos.EnPassant = square
	return nil
}
