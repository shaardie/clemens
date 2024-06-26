package uci

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/shaardie/clemens/pkg/metadata"
	"github.com/shaardie/clemens/pkg/uci/game"
)

func handleInput(input string) {
	tokens := prepareInput(input)
	if len(tokens) == 0 {
		return
	}
	baseCmd := tokens[0]
	tokens = tokens[1:]
	switch baseCmd {
	case "uci":
		fmt.Printf("id name %v %v\nid author %v\nuciok\n", metadata.Name, metadata.Version, metadata.Author)
		return
	case "quit":
		os.Exit(0)
	case "isready":
		g.IsReady()
		return
	case "ucinewgame":
		g = game.New()
		// transpositiontable.Reset()
		return
	case "position":
		g.NewPosition(tokens)
	case "go":
		go g.StartSearch(tokens)
	case "stop":
		g.StopSearch()
	}
}

func prepareInput(s string) []string {
	ss := strings.Fields(s)
	return removePrefixGarbage(ss)
}

var validFirstInputToken = []string{
	"uci", "debug", "isready", "setoption", "ucinewgame", "position", "go", "stop", "quit", "ponderhit",
}

func removePrefixGarbage(ss []string) []string {
	if len(ss) == 0 {
		return ss
	}
	if !slices.Contains(validFirstInputToken, ss[0]) {
		ss = removePrefixGarbage(ss[1:])
	}
	return ss
}
