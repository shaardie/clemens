package main

import (
	"fmt"
	"os"

	"github.com/shaardie/clemens/pkg/uci"
)

func main() {
	err := uci.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
