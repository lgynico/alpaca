package main

import (
	"fmt"
	"os"

	"github.com/lgynico/alpaca/command"
)

func main() {
	if err := command.Run(os.Args); err != nil {
		fmt.Printf("\033[31m[ERROR] %s\033[0m\r\n", err)
	}
}
