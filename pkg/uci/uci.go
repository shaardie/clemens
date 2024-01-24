package uci

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shaardie/clemens/pkg/metadata"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
)

var pos *position.Position

// This is not good
func run(currPos position.Position) {
	m := search.Search(&currPos, 6)
	fmt.Printf("bestmove %v\n", m)
}

func handleInput(input string) string {
	tokens := strings.Fields(input)
	switch tokens[0] {
	case "uci":
		return fmt.Sprintf("id name %v %v\nid author %v\nuciok\n", metadata.Name, metadata.Version, metadata.Author)
	case "quit":
		os.Exit(0)
	case "isready":
		return fmt.Sprintln("readyok")
	case "ucinewgame":
		pos = position.New()
	case "position":
		newPos := handlePositionInput(tokens[1:])
		if newPos != nil {
			pos = newPos
		}
		return ""
	case "go":
		if pos != nil {
			run(*pos)
		}
	}
	return ""
}

func Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimSpace(input)
		output := handleInput(input)
		fmt.Print(output)
	}
	return nil
}
