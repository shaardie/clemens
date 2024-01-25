package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
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
	flag.Parse()
}

func main() {
	pos, err := position.NewFromFen(startPos)
	if err != nil {
		fmt.Printf("Unable to parse fen string, %v", err)
		os.Exit(1)
	}

	before := time.Now()
	results := search.Divided(pos, depth)
	if fen {
		fmt.Println(fenString(results))
	}
	after := time.Since(before)

	if divide {
		fmt.Println(dividedString(results))
	}

	leafs := 0
	for _, r := range results {
		leafs += r.Leafs
	}
	fmt.Printf("Position: %v\nDepth: %d\nLeafs: %d\nDuration: %v\n", startPos, depth, leafs, after)
}

func dividedString(perfts []search.PerftResults) string {
	r := make([]string, len(perfts))
	for i, p := range perfts {
		r[i] = fmt.Sprintf("%v: %d", &p.Move, p.Leafs)
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}

func fenString(perfts []search.PerftResults) string {
	r := make([]string, len(perfts))
	for i, p := range perfts {
		r[i] = fmt.Sprintf("%v %s", &p.Move, p.Position.ToFen())
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}
