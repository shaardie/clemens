package uci

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/shaardie/clemens/pkg/metadata"
	"github.com/shaardie/clemens/pkg/position"
)

func handleInput(input string) {
	tokens := prepareInput(input)
	if len(tokens) == 0 {
		return
	}
	switch tokens[0] {
	case "uci":
		fmt.Printf("id name %v %v\nid author %v\nuciok\n", metadata.Name, metadata.Version, metadata.Author)
		return
	case "quit":
		os.Exit(0)
	case "isready":
		fmt.Println("readyok")
		return
	case "ucinewgame":
		pos = position.New()
		return
	case "position":
		newPos := handlePositionInput(tokens[1:])
		if newPos != nil {
			pos = newPos
		}
	case "go":
		if pos != nil {
			run(*pos)
		}
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
