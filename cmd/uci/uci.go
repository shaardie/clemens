package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
)

var pos *position.Position

// This is not good
func run(currPos position.Position) {
	m := search.Search(&currPos, 5)
	fmt.Printf("bestmove %v\n", m)
}

func handleInput(input string) string {
	tokens := strings.Fields(input)
	switch tokens[0] {
	case "uci":
		return "id name Clemens 0.1.0\nid author Sven Haardiek <sven@haardiek.de>\nuciok\n"
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

func mainWithError() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimSpace(input)
		output := handleInput(input)
		fmt.Print(output)
	}
	return nil
}

func main() {
	err := mainWithError()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
