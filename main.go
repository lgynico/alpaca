package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/lgynico/alpaca/mate"
)

var (
	config_dir string
	json_out   string
	go_out     string

	w = flag.CommandLine.Output()
)

func init() {
	flag.StringVar(&config_dir, "dir", "", "path to excel config files")
	flag.StringVar(&json_out, "json_out", "", "path to output json files")
	flag.StringVar(&go_out, "go_out", "", "path to output golang files")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	if err := checkFlags(); err != nil {
		fmt.Fprintln(w, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	mates, err := mate.Parse(config_dir)
	if err != nil {
		panic(err)
	}

	for _, exec := range executors {
		if err = exec(mates); err != nil {
			panic(err)
		}
	}
}

func usage() {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  alpaca [flags]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Flags:")
	flag.CommandLine.PrintDefaults()
}

func checkFlags() error {
	if len(config_dir) == 0 {
		return errors.New("flag -dir is require")
	}

	if len(json_out) == 0 && len(go_out) == 0 {
		return errors.New("specify at least one *_out flag")
	}

	return nil
}
