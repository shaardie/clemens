package main

import (
	"fmt"
	"os"

	"github.com/shaardie/clemens/pkg/uci"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("info '%v'", err)
		}
	}()
	err := uci.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
