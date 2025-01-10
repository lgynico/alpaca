package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/types"
)

var (
	config_dir string
	json_out   string
	go_out     string
	cs_out     string

	w = flag.CommandLine.Output()
)

var (
	configMetas []*meta.Config
)

func init() {
	flag.StringVar(&config_dir, "dir", "", "path to excel config files")
	flag.StringVar(&json_out, "json_out", "", "path to output json files")
	flag.StringVar(&go_out, "go_out", "", "path to output golang files")
	flag.StringVar(&cs_out, "cs_out", "", "path to output c# files")

	flag.Usage = usage
}

func main() {
	checkFlag()
	parseEnum()
	parseConfig()
	generateFiles()
}

func parseEnum() {
	if err := types.ParseEnum(config_dir); err != nil {
		panic(err)
	}
}

func parseConfig() {
	metas, err := meta.Parse(config_dir)
	if err != nil {
		panic(err)
	}

	for _, exec := range executors {
		if err = exec(metas); err != nil {
			panic(err)
		}
	}

	configMetas = metas
}

func checkFlag() {
	flag.Parse()

	if err := checkFlags(); err != nil {
		_, _ = fmt.Fprintln(w, err.Error())
		flag.Usage()
		os.Exit(1)
	}
}

func usage() {
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  alpaca [flags]")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Flags:")
	flag.CommandLine.PrintDefaults()
}

func checkFlags() error {
	if len(config_dir) == 0 {
		return errors.New("flag -dir is require")
	}

	if len(json_out) == 0 && len(go_out) == 0 && len(cs_out) == 0 {
		return errors.New("specify at least one *_out flag")
	}

	return nil
}

func generateFiles() {
	for _, f := range generators {
		if err := f(configMetas); err != nil {
			panic(err)
		}
	}
}
