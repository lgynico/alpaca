package value

import (
	"fmt"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/mate"
)

var Parser = &parser{}

type parser struct {
}

func (p *parser) Visit(configMate *mate.Config) error {
	for _, field := range configMate.Fields {
		if err := p.parseFieldValues(field); err != nil {
			return fmt.Errorf("value parse error on field %s: %v", field.Name, err)
		}
	}
	return nil
}

func (p *parser) parseFieldValues(f *mate.Field) error {
	f.Values = make([]any, len(f.RawValues))
	for i, rawValue := range f.RawValues {
		value, err := consts.ParseValue(rawValue, f.Type, f.TypeParams...)
		if err != nil {
			return err
		}

		f.Values[i] = value
	}
	return nil
}
