package main

import (
	"fmt"

	"github.com/shaardie/clemens/pkg/position"
)

func main() {
	fmt.Println(position.New().ToFen())

}
