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
	zobristHash uint64
	bestMove    move.Move
	score       int16
	depth       uint8

	// 0-1 for the Node Type
	// 2-7 for the Age
	ageAndNodeType uint8
}

func (te *TranspositionEntry) getNodeType() nodeType {
	return nodeType(te.ageAndNodeType & 0b11)
}
func (te *TranspositionEntry) setNodeType(nt nodeType) {
	te.ageAndNodeType = te.ageAndNodeType&0b11111100 | uint8(nt)
}
func (te *TranspositionEntry) getAge() uint8 {
	return te.ageAndNodeType >> 2
}
func (te *TranspositionEntry) setAge(age uint8) {
	te.ageAndNodeType = te.ageAndNodeType&0b11 | (age << 2)
}

// 64MB
const transpositionTableSizeinMB = 1024 * 1024 * 128

var transpositionTableSize uint64
var (
	tt []TranspositionEntry
)

var hashEntries uint64

func HashFull() uint64 {
	return 1000 * hashEntries / transpositionTableSize
}

func init() {
	Reset()
}

func Reset() {
	transpositionTableSize = uint64(transpositionTableSizeinMB / unsafe.Sizeof(TranspositionEntry{}))
	tt = make([]TranspositionEntry, transpositionTableSize)
}

func Get(zobristHash uint64, alpha, beta int16, depth, ply uint8) (score int16, use bool, m move.Move) {
	key := zobristHash % transpositionTableSize
	te := &tt[key]

	if te.zobristHash != zobristHash {
		return 0, false, move.NullMove
	}

	// Only use the value, if the depth of the entry is bigger that the current one.
	// Remember that the depth decreases, while going down the tree.
	// If we use this entry or not. The move should be a good guess.
	if te.depth < depth {
		return 0, false, te.bestMove
	}

	score = te.score

	// // Adjust if mate value
	if score > evaluation.INF-100 {
		score -= int16(ply)
	} else if score < -evaluation.INF+100 {
		score += int16(ply)
	}

	switch te.getNodeType() {
	case AlphaNode:
		if score <= alpha {
			return alpha, true, te.bestMove
		}
	case BetaNode:
		if score >= beta {
			return beta, true, te.bestMove
		}
	}

	return score, true, te.bestMove
}

// PotentiallySave save the new transposition entry, if it is a better fit.
// Note, that we use single values as parameter for the case, so we not create the struct, if we do not have to
func PotentiallySave(zobristHash uint64, bestMove move.Move, depth uint8, score int16, nt nodeType, age uint8) {
	key := zobristHash % transpositionTableSize
	te := &tt[key]

	// Ignore the new entry, if there is an entry with a higher depth.
	if te.depth > depth || te.getAge() < age {
		return
	}

	// New Entry
	if te.zobristHash == 0 {
		hashEntries++
	}

	te.zobristHash = zobristHash
	te.bestMove = bestMove
	te.depth = depth
	te.score = score
	te.setNodeType(nt)
	te.setAge(age)
}
