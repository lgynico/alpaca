package rule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/mate"
)

var Checker = &checker{
	ruleCheckers: map[consts.RuleType]checkFunc{
		consts.NoRule:      noRuleCheck,
		consts.KeyRule:     keyRuleCheck,
		consts.UniqueRule:  uniqueRuleCheck,
		consts.RequireRule: requireRuleCheck,
		consts.RangeRule:   rangeRuleCheck,
		consts.LengthRule:  lengthRuleCheck,
		consts.DecimalRule: decimalRuleCheck,
	},
}

type (
	checkFunc func(meta *mate.Config, field *mate.Field) error

	checker struct {
		ruleCheckers map[consts.RuleType]checkFunc
	}
)

func (c *checker) Visit(configMeta *mate.Config) error {
	for _, field := range configMeta.Fields {
		checkFunc := c.getCheckFunc(consts.RuleType(field.RuleMeta.Key))
		if err := checkFunc(configMeta, field); err != nil {
			return fmt.Errorf("rule check error on field %s: %v", field.Name, err)
		}
	}
	return nil
}

func (c *checker) getCheckFunc(ruleType consts.RuleType) checkFunc {
	if ruleChecker, ok := c.ruleCheckers[ruleType]; ok {
		return ruleChecker
	}

	return errRuleCheck
}

func noRuleCheck(meta *mate.Config, field *mate.Field) error {
	return nil
}

func errRuleCheck(meta *mate.Config, field *mate.Field) error {
	return fmt.Errorf("unknown rule string: %s", field.RuleMeta.Origin)
}

func keyRuleCheck(meta *mate.Config, field *mate.Field) error {
	if meta.KeyField != nil {
		return errors.New("duplicate key field")
	}

	meta.KeyField = field
	return uniqueRuleCheck(meta, field)
}

func uniqueRuleCheck(meta *mate.Config, field *mate.Field) error {
	set := map[string]bool{}
	for _, value := range field.RawValues {
		if _, ok := set[value]; ok {
			return fmt.Errorf("duplicate value on unique field: %s, value: %s", field.Name, value)
		}

		set[value] = true
	}
	return nil
}

func requireRuleCheck(meta *mate.Config, field *mate.Field) error {
	for _, value := range field.RawValues {
		if len(strings.TrimSpace(value)) == 0 {
			return fmt.Errorf("empty value on require field: %s", field.Name)
		}
	}
	return nil
}

func rangeRuleCheck(meta *mate.Config, field *mate.Field) error {
	for _, value := range field.RawValues {
		switch field.Type {
		case consts.Int, consts.Int8, consts.Int32, consts.Int64:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}

			if err = checkRange[int64](i, field.RuleMeta.Params, parseInt()); err != nil {
				return err
			}

		case consts.Uint, consts.Uint8, consts.Uint32, consts.Uint64:
			i, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}

			if err = checkRange[uint64](i, field.RuleMeta.Params, parseUint()); err != nil {
				return err
			}

		case consts.Float, consts.Double:
			i, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}

			if err = checkRange[float64](i, field.RuleMeta.Params, parseFloat()); err != nil {
				return err
			}

		default:
			return fmt.Errorf("range rule only supports for numeric types, not for %s type", field.Type)
		}
	}

	return nil
}

func lengthRuleCheck(meta *mate.Config, field *mate.Field) error {
	if field.Type != consts.String {
		return fmt.Errorf("range rule only supports for string types, not for %s type", field.Type)
	}

	min, err := strconv.ParseUint(field.RuleMeta.Params[0], 10, 64)
	if err != nil {
		return err
	}
	max, err := strconv.ParseUint(field.RuleMeta.Params[1], 10, 64)
	if err != nil {
		return err
	}

	for _, value := range field.RawValues {
		length := uint64(len(value))
		if length < min || length > max {
			return errors.New("string length out of range")
		}
	}

	return nil
}

func decimalRuleCheck(meta *mate.Config, field *mate.Field) error {
	// TODO: trim decimal
	return nil
}

func checkRange[E ~int64 | ~uint64 | ~float64](i E, params []string, parser func(value string) (E, error)) error {
	min, err := parser(params[1])
	if err != nil {
		return err
	}

	max, err := parser(params[2])
	if err != nil {
		return err
	}

	outOfRange := (params[0] == "(" && i <= min) ||
		(params[0] == "[" && i < min) ||
		(params[3] == ")" && i >= max) ||
		(params[3] == "]" && i > max)

	if outOfRange {
		return errors.New("value out of range")
	}

	return nil
}

func parseInt() func(value string) (int64, error) {
	return func(value string) (int64, error) {
		return strconv.ParseInt(value, 10, 64)
	}
}

func parseUint() func(value string) (uint64, error) {
	return func(value string) (uint64, error) {
		return strconv.ParseUint(value, 10, 64)
	}
}

func parseFloat() func(value string) (float64, error) {
	return func(value string) (float64, error) {
		return strconv.ParseFloat(value, 64)
	}
}
