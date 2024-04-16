package evaluation

import (
	"fmt"
	"strings"

	"github.com/shaardie/clemens/pkg/types"
)

const NUMBER_OF_WEIGHTS = len(PieceValue) - 1 +
	2 + // king shield value
	len(kingAttValue) + // kind attack values
	3 + // Pairs
	int(types.PIECE_TYPE_NUMBER)*int(game_number)*int(types.SQUARE_NUMBER) + // Piece Square Table
	len(knight_pawn_adjustment) + len(rook_pawn_adjustment) + // Pawn Adjustments
	len(isolanis) + len(doubled) + len(passPawns) + len(backwards) + len(supported) + len(phalanx) + len(opposed) + // Pawn Evaluation
	4 // Rook Evaluation

type Weights [NUMBER_OF_WEIGHTS]int16

func prettify(is []int16) string {
	builder := strings.Builder{}
	builder.WriteRune('{')
	for i := range len(is) {
		if i%8 == 0 {
			builder.WriteString("\n\t")
		}
		builder.WriteString(fmt.Sprintf("%v, ", is[i]))
	}
	builder.WriteString("\n")
	builder.WriteRune('}')
	return builder.String()
}

func PrintWeights(weights Weights) {
	idx := 0
	fmt.Printf("PieceValue: %v\n", prettify(weights[idx:idx+len(PieceValue)-1]))
	idx += len(PieceValue) - 1

	fmt.Printf("King Shield 2: %v\n", weights[idx])
	idx++
	fmt.Printf("King Shield 3: %v\n", weights[idx])
	idx++

	fmt.Printf("King Attack Values: %v\n", prettify(weights[idx:idx+len(kingAttValue)]))
	idx += len(kingAttValue)

	fmt.Printf("Rook Pair: %v\n", weights[idx])
	idx++
	fmt.Printf("Knight Pair: %v\n", weights[idx])
	idx++
	fmt.Printf("Bishop Pair: %v\n", weights[idx])
	idx++

	for pieceType := range len(pieceTables) {
		for phase := range len(pieceTables[pieceType]) {
			fmt.Printf("Piece Table '%v' '%v': %v\n", types.PieceType(pieceType), gamePhaseString(phase), prettify(weights[idx:idx+len(pieceTables[pieceType][phase])]))
			idx += len(pieceTables[pieceType][phase])
		}
	}
	fmt.Println(idx)
	fmt.Printf("Knight Pawn Adjustments: %v\n", prettify(weights[idx:idx+len(knight_pawn_adjustment)]))
	idx += len(knight_pawn_adjustment)

	fmt.Printf("Rook Pawn Adjustments: %v\n", prettify(weights[idx:idx+len(rook_pawn_adjustment)]))
	idx += len(rook_pawn_adjustment)

	fmt.Printf("Isolanis: %v\n", prettify(weights[idx:idx+len(isolanis)]))
	idx += len(isolanis)

	fmt.Printf("Doubled Pawns: %v\n", prettify(weights[idx:idx+len(doubled)]))
	idx += len(doubled)

	fmt.Printf("Passed Pawns: %v\n", prettify(weights[idx:idx+len(passPawns)]))
	idx += len(passPawns)

	fmt.Printf("Backwards Pawns: %v\n", prettify(weights[idx:idx+len(backwards)]))
	idx += len(backwards)
	fmt.Printf("Supported Pawns: %v\n", prettify(weights[idx:idx+len(supported)]))
	idx += len(supported)
	fmt.Printf("Phalanxed Pawns: %v\n", prettify(weights[idx:idx+len(phalanx)]))
	idx += len(phalanx)
	fmt.Printf("Opposed Pawns: %v\n", prettify(weights[idx:idx+len(opposed)]))
	idx += len(opposed)

	fmt.Printf("Rook Seventh Midgame: %v\n", weights[idx])
	idx++
	fmt.Printf("Rook Seventh Endgame: %v\n", weights[idx])
	idx++
	fmt.Printf("Rook Open File: %v\n", weights[idx])
	idx++
	fmt.Printf("Rook Halfopen File: %v\n", weights[idx])
	idx++
}

func SetFromWeights(weights Weights) {
	idx := 0
	for i := range len(PieceValue) - 1 {
		PieceValue[i] = weights[idx]
		idx++
	}

	shield2Value = weights[idx]
	idx++
	shield3Value = weights[idx]
	idx++

	for i := range len(kingAttValue) {
		kingAttValue[i] = weights[idx]
		idx++
	}

	rookPair = weights[idx]
	idx++
	knightPair = weights[idx]
	idx++
	knightPair = weights[idx]
	idx++

	for pieceType := range len(pieceTables) {
		for game := range len(pieceTables[pieceType]) {
			for square := range len(pieceTables[pieceType][game]) {
				pieceTables[pieceType][game][square] = weights[idx]
				idx++
			}
		}
	}
	initPieceSquareTables()

	for i := range len(knight_pawn_adjustment) {
		knight_pawn_adjustment[i] = weights[idx]
		idx++
	}
	for i := range len(rook_pawn_adjustment) {
		rook_pawn_adjustment[i] = weights[idx]
		idx++
	}

	for i := range len(isolanis) {
		isolanis[i] = weights[idx]
		idx++
	}
	for i := range len(doubled) {
		doubled[i] = weights[idx]
		idx++
	}
	for i := range len(passPawns) {
		passPawns[i] = weights[idx]
		idx++
	}
	for i := range len(backwards) {
		backwards[i] = weights[idx]
		idx++
	}
	for i := range len(supported) {
		supported[i] = weights[idx]
		idx++
	}
	for i := range len(phalanx) {
		phalanx[i] = weights[idx]
		idx++
	}
	for i := range len(opposed) {
		opposed[i] = weights[idx]
		idx++
	}

	rookSeventhMidgame = weights[idx]
	idx++
	rookSeventhEndgame = weights[idx]
	idx++
	rookOpenFile = weights[idx]
	idx++
	rookHalfOpenFile = weights[idx]
	idx++
}

func GetWeights() Weights {
	var weights Weights
	idx := 0
	for i := range len(PieceValue) - 1 {
		weights[idx] = PieceValue[i]
		idx++
	}

	weights[idx] = shield2Value
	idx++
	weights[idx] = shield2Value
	idx++

	for i := range len(kingAttValue) {
		weights[idx] = kingAttValue[i]
		idx++
	}

	weights[idx] = rookPair
	idx++
	weights[idx] = knightPair
	idx++
	weights[idx] = knightPair
	idx++

	for pieceType := range len(pieceTables) {
		for game := range len(pieceTables[pieceType]) {
			for square := range len(pieceTables[pieceType][game]) {
				weights[idx] = pieceTables[pieceType][game][square]
				idx++
			}
		}
	}

	for i := range len(knight_pawn_adjustment) {
		weights[idx] = knight_pawn_adjustment[i]
		idx++
	}

	for i := range len(rook_pawn_adjustment) {
		weights[idx] = rook_pawn_adjustment[i]
		idx++
	}

	for i := range len(isolanis) {
		weights[idx] = isolanis[i]
		idx++
	}
	for i := range len(doubled) {
		weights[idx] = doubled[i]
		idx++
	}
	for i := range len(passPawns) {
		weights[idx] = passPawns[i]
		idx++
	}
	for i := range len(backwards) {
		weights[idx] = backwards[i]
		idx++
	}
	for i := range len(supported) {
		weights[idx] = supported[i]
		idx++
	}
	for i := range len(phalanx) {
		weights[idx] = phalanx[i]
		idx++
	}
	for i := range len(opposed) {
		weights[idx] = opposed[i]
		idx++
	}

	weights[idx] = rookSeventhMidgame
	idx++
	weights[idx] = rookSeventhEndgame
	idx++
	weights[idx] = rookOpenFile
	idx++
	weights[idx] = rookHalfOpenFile
	idx++

	return weights
}
