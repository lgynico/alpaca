package main

import (
	"testing"

	"github.com/lgynico/alpaca/mate"
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
