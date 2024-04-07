package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

var (
	startPos string
	depth    int
	divide   bool
	fen      bool
)

func init() {
	flag.StringVar(&startPos, "position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "start position for perft test")
	flag.IntVar(&depth, "depth", 1, "depth for perft test")
	flag.BoolVar(&divide, "divide", false, "print divided output")
	flag.BoolVar(&fen, "fen", false, "print fen strings for the positions in the first depth")

}

func main() {
	flag.Parse()
	pos, err := position.NewFromFen(startPos)
	if err != nil {
		fmt.Printf("Unable to parse fen string, %v", err)
		os.Exit(1)
	}

	var duration time.Duration
	var leafs int

	if divide {
		before := time.Now()
		results := Divided(pos, depth)
		duration = time.Since(before)
		if fen {
			fmt.Println(fenString(results))
		}
		fmt.Println(dividedString(results))
		for _, r := range results {
			leafs += r.Leafs
		}
	} else {
		before := time.Now()
		leafs = Perft(pos, depth)
		duration = time.Since(before)
	}

	fmt.Printf("Position: %v\nDepth: %d\nLeafs: %d\nDuration: %v\nNodes per Second: %v\n",
		startPos,
		depth,
		leafs,
		duration,
		PerftNodes*1000/(max(duration.Milliseconds(), 1)),
	)
}

func dividedString(perfts []PerftResults) string {
	r := make([]string, len(perfts))
	for i, p := range perfts {
		r[i] = fmt.Sprintf("%v: %d", &p.Move, p.Leafs)
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}

func fenString(perfts []PerftResults) string {
	r := make([]string, len(perfts))
	for i, p := range perfts {
		r[i] = fmt.Sprintf("%v %s", &p.Move, p.Position.ToFen())
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}

type PerftResults struct {
	Move     move.Move
	Position position.Position
	Leafs    int
}

var PerftNodes int64 = 0

func Perft(pos *position.Position, depth int) int {
	PerftNodes++
	if depth == 0 {
		return 1
	}
	var leafs int
	var prevPos position.Position
	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(*m)
		if pos.IsLegal() {
			leafs += Perft(pos, depth-1)
		}
		*pos = prevPos
	}
	return leafs
}

func Divided(pos *position.Position, depth int) []PerftResults {
	PerftNodes++
	if depth == 0 {
		panic("depth should be bigger than 0")
	}

	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	results := make([]PerftResults, 0, moves.Length())
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos := *pos
		pos.MakeMove(*m)
		if pos.IsLegal() {
			results = append(results,
				PerftResults{
					Move:     *m,
					Position: *pos,
					Leafs:    Perft(pos, depth-1),
				},
			)
		}
		pos = &prevPos
	}
	return results
}
