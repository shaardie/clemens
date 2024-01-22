package main

import (
	"strings"

	"github.com/shaardie/clemens/pkg/position"
)

func handlePositionInput(tokens []string) *position.Position {
	var pos *position.Position
	if len(tokens) == 0 {
		return nil
	}
	switch tokens[0] {
	case "startpos":
		pos = position.New()
		tokens = tokens[1:]
	case "fen":
		if len(tokens) < 7 {
			return nil
		}
		fenPos, err := position.NewFromFen(strings.Join(tokens[1:7], " "))
		if err != nil {
			return nil
		}
		pos = fenPos
		tokens = tokens[7:]
	}
	if len(tokens) <= 1 || tokens[0] != "moves" {
		return pos
	}
	tokens = tokens[1:]

	for _, token := range tokens {
		err := pos.MakeMoveFromString(token)
		if err != nil {
			return pos
		}
	}
	return pos
}
