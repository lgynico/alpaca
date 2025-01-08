package main

import (
	"fmt"

	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/rule"
	"github.com/lgynico/alpaca/value"
)

type executor func(metaList []*meta.Config) error

var executors = []executor{
	parseRule,
	parseValue,
	checkRule,
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
			return err
		}
		fmt.Printf("parse [%s] value SUCCEED !\r\n", meta.Filename)
	}

	fmt.Println("< parse values SUCCEED")
	return nil
}
