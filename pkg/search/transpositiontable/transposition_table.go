package transpositiontable

import (
	"unsafe"

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
var TTable TranspositionTable

type TranspositionTable []TranspositionEntry

func init() {
	reset()
}

func reset() {
	transpositionTableSize = uint64(transpositionTableSizeinMB / unsafe.Sizeof(TranspositionEntry{}))
	TTable = make([]TranspositionEntry, transpositionTableSize)
}

func (tt TranspositionTable) Get(zobristHash uint64, depth uint8) (te TranspositionEntry, found bool, isGoodGuess bool) {
	key := zobristHash % transpositionTableSize
	te = tt[key]
	if te.ZobristHash != zobristHash {
		return te, false, false
	}

	// Only use the value, if the depth of the entry is bigger that the current one.
	// Remember that the depth decreases, while going down the tree.
	// If we use this entry or not. The move should be a good guess.
	if te.Depth < depth {
		return te, false, true
	}

	return te, true, true
}

// PotentiallySave save the new transposition entry, if it is a better fit.
// Note, that we use single values as parameter for the case, so we not create the struct, if we do not have to
func (tt TranspositionTable) PotentiallySave(zobristHash uint64, bestMove move.Move, depth uint8, score int, nt nodeType) {
	key := zobristHash % transpositionTableSize
	oldTe := tt[key]

	// Ignore the new entry, if there is a mathing entry with a higher depth
	if oldTe.ZobristHash == zobristHash && oldTe.Depth > depth {
		return
	}

	tt[key] = TranspositionEntry{
		ZobristHash: zobristHash,
		BestMove:    bestMove,
		Depth:       depth,
		Score:       score,
		NodeType:    nt,
	}
}
