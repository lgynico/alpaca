package main

import (
	"fmt"
	"os"

	"github.com/lgynico/alpaca/command"
)

func main() {
	if err := command.Run(os.Args); err != nil {
		fmt.Printf("%v\n", err)
	}
}
