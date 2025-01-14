package command

import (
	"fmt"

	"github.com/lgynico/alpaca/meta"
	"github.com/lgynico/alpaca/rule"
	"github.com/lgynico/alpaca/types"
	"github.com/lgynico/alpaca/value"
)

func (p *app) parseFiles() ([]*meta.Config, error) {
	dir := p.input

	if err := types.ParseEnum(dir); err != nil {
		return nil, err
	}

	metas, err := meta.Parse(dir)
	if err != nil {
		return nil, err
	}

	if err = p.parseRule(metas); err != nil {
		return nil, err
	}

	if err = p.checkRule(metas); err != nil {
		return nil, err
	}

	if err = p.parseValue(metas); err != nil {
		return nil, err
	}

	return metas, nil
}

func (p *app) parseRule(metas []*meta.Config) error {
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

func (p *app) checkRule(metas []*meta.Config) error {
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
func (p *app) parseValue(metas []*meta.Config) error {
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
