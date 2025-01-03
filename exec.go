package main

import (
	"log"
	"os/exec"

	"github.com/lgynico/alpaca/mate"
	"github.com/lgynico/alpaca/rule"
	"github.com/lgynico/alpaca/value"
	"github.com/lgynico/alpaca/writer"
)

type executor func(metaList []*mate.Config) error

var executors = []executor{
	parseRule,
	parseValue,
	checkRule,

	genJSON,
	genGO,
}

func parseRule(mates []*mate.Config) error {
	log.Println("parse rules ...")
	for _, mate := range mates {
		if err := rule.Parser.Visit(mate); err != nil {
			return err
		}
		log.Printf("parse [%s] rule SUCCEED !", mate.Filename)
	}

	log.Println("parse rules SUCCEED")
	return nil
}

func checkRule(mates []*mate.Config) error {
	log.Println("check rules ...")
	for _, mate := range mates {
		if err := rule.Checker.Visit(mate); err != nil {
			return err
		}
		log.Printf("check [%s] rule SUCCEED !", mate.Filename)
	}

	log.Println("check rules SUCCEED")
	return nil
}
func parseValue(mates []*mate.Config) error {
	log.Println("parse values ...")
	for _, mate := range mates {
		if err := value.Parser.Visit(mate); err != nil {
			return nil
		}
		log.Printf("parse [%s] value SUCCEED !", mate.Filename)
	}

	log.Println("parse values SUCCEED")
	return nil
}

func genJSON(mates []*mate.Config) error {
	if len(json_out) == 0 {
		return nil
	}

	for _, mate := range mates {
		if err := writer.WriteJSON(json_out, mate); err != nil {
			return err
		}
	}

	return nil
}

func genGO(mates []*mate.Config) error {
	if len(go_out) == 0 {
		return nil
	}

	for _, mate := range mates {
		if err := writer.WriteGoConfig(go_out, mate); err != nil {
			return err
		}
	}

	if err := writer.WriteGoConfigMgr(go_out, mates); err != nil {
		return err
	}

	return exec.Command("gofmt", "-w", go_out).Run()
}
