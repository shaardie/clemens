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

// 64MB
const ttSizeInMB = 1024 * 1024 * 64

type bucket []ttEntry

const numberOfBuckets = 4

var bucketSize uint64

var tt = [numberOfBuckets]bucket{}

var hashEntries uint64

func HashFull() uint64 {
	return 1000 * hashEntries / (bucketSize * numberOfBuckets)
}

func init() {
	Reset()
}

func Reset() {
	bucketSize = ttSizeInMB / numberOfBuckets / uint64(unsafe.Sizeof(ttEntry{}))
	for i := range tt {
		if tt[i] == nil {
			tt[i] = make(bucket, bucketSize)
		} else {
			clear(tt[i])
		}
	}
}

func Get(zobristHash uint64, alpha, beta int16, depth, ply uint8) (score int16, use bool, m move.Move) {
	key := zobristHash % bucketSize
	var te *ttEntry
	var found bool
	for i := range tt {
		te = &tt[i][key]
		if te.zobristHash == zobristHash {
			found = true
			break
		}
	}

	// No entry found
	if !found {
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
	case PVNode:
		return score, true, te.bestMove
	}

	return score, false, te.bestMove
}

// PotentiallySave save the new transposition entry, if it is a better fit.
// Note, that we use single values as parameter for the case, so we not create the struct, if we do not have to
func PotentiallySave(zobristHash uint64, bestMove move.Move, depth uint8, score int16, nt nodeType, age uint8) {
	var te *ttEntry
	key := zobristHash % bucketSize
	for i := range tt {
		te = &tt[i][key]

		// Empty Entries should always be overriden
		if te.zobristHash == 0 {
			hashEntries++
			break
		}

		// Found a worse one, replace
		if te.depth <= depth && te.getAge() >= age {
			break
		}

		// If none is found, we replace the last one
	}

	te.zobristHash = zobristHash
	te.bestMove = bestMove
	te.depth = depth
	te.score = score
	te.setNodeType(nt)
	te.setAge(age)
}
