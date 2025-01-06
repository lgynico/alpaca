package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lgynico/alpaca/types"

	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/rule"
	"github.com/lgynico/alpaca/value"
	"github.com/lgynico/alpaca/writer"
)

type executor func(metaList []*meta.Config) error

var executors = []executor{
	parseRule,
	parseValue,
	checkRule,

	genJSON,
	genGO,
}

func parseRule(metas []*meta.Config) error {
	fmt.Println("> parse rules ...")
	for _, meta := range metas {
		if err := rule.Parser.Visit(meta); err != nil {
			return err
		}
		fmt.Printf("parse [%s] rule SUCCEED !\r\n", meta.Filename)
	}

	fmt.Println("< parse rules SUCCEED")
	return nil
}

func checkRule(metas []*meta.Config) error {
	fmt.Println("> check rules ...")
	for _, meta := range metas {
		if err := rule.Checker.Visit(meta); err != nil {
			return err
		}
		fmt.Printf("check [%s] rule SUCCEED !\r\n", meta.Filename)
	}

	fmt.Println("< check rules SUCCEED")
	return nil
}
func parseValue(metas []*meta.Config) error {
	fmt.Println("> parse values ...")
	for _, meta := range metas {
		if err := value.Parser.Visit(meta); err != nil {
			return nil
		}
		fmt.Printf("parse [%s] value SUCCEED !\r\n", meta.Filename)
	}

	fmt.Println("< parse values SUCCEED")
	return nil
}

func genJSON(metas []*meta.Config) error {
	if len(json_out) == 0 {
		return nil
	}

	mkdir(json_out)

	fmt.Println("> write json ...")
	for _, meta := range metas {
		if err := writer.WriteJSON(json_out, meta); err != nil {
			return err
		}
		fmt.Printf("write [%s.json] SUCCEED !\r\n", meta.Filename)
	}
	fmt.Println("< write json SUCCEED !")

	return nil
}

func genGO(metas []*meta.Config) error {
	if len(go_out) == 0 {
		return nil
	}

	mkdir(go_out)

	fmt.Println("> write go configs ...")
	if err := writer.WriteGoConfigs(go_out, metas); err != nil {
		return err
	}
	fmt.Println("< write go configs SUCCEED !")

	if err := writer.WriteGoConfigMgr(go_out, metas); err != nil {
		return err
	}
	fmt.Println("< write go ConfigMgr SUCCEED !")

	if err := writer.WriteGoEnums(go_out, types.Enums()); err != nil {
		return err
	}
	fmt.Println("< write go enums SUCCEED !")

	if err := os.Chdir(go_out); err == nil {
		if err := exec.Command("gofmt", "-w", ".").Run(); err != nil {
			fmt.Printf("format go codes FAILED: %v\r\n", err)
		} else {
			fmt.Println("< format go codes SUCCEED !")
		}
	} else {
		fmt.Printf("format go codes FAILED: %v\r\n", err)
	}

	return nil
}

func mkdir(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
}
