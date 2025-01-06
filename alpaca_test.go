package main

import (
	"fmt"
	"testing"

	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/types"
)

func TestCheckRule(t *testing.T) {

	config_dir = "./example"
	go_out = "./example"

	metas, err := meta.Parse(config_dir)
	if err != nil {
		panic(err)
	}

	if err = parseRule(metas); err != nil {
		t.Fatal(err)
	}
	if err = parseValue(metas); err != nil {
		t.Fatal(err)
	}
	if err = checkRule(metas); err != nil {
		t.Fatal(err)
	}

}

func TestEnums(t *testing.T) {
	config_dir = "./example"
	if err := types.ParseEnum(config_dir); err != nil {
		t.Fatal(err)
	}

	for name, enum := range types.Enums() {
		fmt.Println("enum:", name)
		for _, node := range enum.Nodes {
			fmt.Printf("%s=%d    ", node.Key, node.Value)
		}
		fmt.Println()
	}
}
