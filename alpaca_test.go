package main

import (
	"fmt"
	"testing"

	"github.com/lgynico/alpaca/mate"
	"github.com/lgynico/alpaca/types"
)

func TestCheckRule(t *testing.T) {

	config_dir = "./example"
	go_out = "./example"

	mates, err := mate.Parse(config_dir)
	if err != nil {
		panic(err)
	}

	if err = parseRule(mates); err != nil {
		t.Fatal(err)
	}
	if err = parseValue(mates); err != nil {
		t.Fatal(err)
	}
	if err = checkRule(mates); err != nil {
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
