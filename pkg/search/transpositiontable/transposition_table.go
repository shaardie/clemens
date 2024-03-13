package transpositiontable

import (
	"unsafe"

	"github.com/shaardie/clemens/pkg/evaluation"
	"github.com/shaardie/clemens/pkg/move"
)

type nodeType uint8

const (
	PVNode nodeType = iota
	AlphaNode
	BetaNode
)

type TranspositionEntry struct {
	ZobristHash uint64
	BestMove    move.Move
	Depth       uint8
	Score       int
	NodeType    nodeType
}

// 64MB
const transpositionTableSizeinMB = 1024 * 1024 * 64

var transpositionTableSize uint64
var (
	tt []TranspositionEntry
)

var hashEntries uint64

func HashFull() uint64 {
	return 1000 * hashEntries / transpositionTableSize
}

func init() {
	reset()
}

func reset() {
	transpositionTableSize = uint64(transpositionTableSizeinMB / unsafe.Sizeof(TranspositionEntry{}))
	tt = make([]TranspositionEntry, transpositionTableSize)
}

func Get(zobristHash uint64, alpha, beta int, depth, ply uint8) (score int, use bool, m move.Move) {
	key := zobristHash % transpositionTableSize
	te := &tt[key]

	if te.ZobristHash != zobristHash {
		return 0, false, move.NullMove
	}

	// Only use the value, if the depth of the entry is bigger that the current one.
	// Remember that the depth decreases, while going down the tree.
	// If we use this entry or not. The move should be a good guess.
	if te.Depth < depth {
		return 0, false, te.BestMove
	}

	score = te.Score

	// // Adjust if mate value
	if score > evaluation.INF-100 {
		score -= int(ply)
	} else if score < -evaluation.INF+100 {
		score += int(ply)
	}

	switch te.NodeType {
	case AlphaNode:
		if score <= alpha {
			return alpha, true, te.BestMove
		}
	case BetaNode:
		if score >= beta {
			return beta, true, te.BestMove
		}
	}

	return score, true, te.BestMove
}

// PotentiallySave save the new transposition entry, if it is a better fit.
// Note, that we use single values as parameter for the case, so we not create the struct, if we do not have to
func PotentiallySave(zobristHash uint64, bestMove move.Move, depth uint8, score int, nt nodeType) {
	key := zobristHash % transpositionTableSize
	oldTe := tt[key]

	// Ignore the new entry, if there is an entry with a higher depth.
	if oldTe.Depth > depth {
		return
	}

	// New Entry
	if oldTe.ZobristHash == 0 {
		hashEntries++
	}

	tt[key] = TranspositionEntry{
		ZobristHash: zobristHash,
		BestMove:    bestMove,
		Depth:       depth,
		Score:       score,
		NodeType:    nt,
	}
}
