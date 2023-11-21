package rule

import (
	"regexp"

	"github.com/lgynico/alpaca/mate"
)

var Parser = &parser{
	patterns: []*regexp.Regexp{
		regexp.MustCompile(`^(key)$`),
		regexp.MustCompile(`^(unique)$`),
		regexp.MustCompile(`^(require)$`),
		regexp.MustCompile(`^(range)\s*([\(\[])\s*(\d+)\s*,\s*(\d+)\s*([\)\]])$`),
		regexp.MustCompile(`^(length)\s*\[\s*(\d+)\s*,\s*(\d+)\s*\]$`),
		regexp.MustCompile(`^(decimal)\s*:\s*(\d+)$`),
	},
}

type parser struct {
	patterns []*regexp.Regexp
}

func (p *parser) Visit(configMeta *mate.Config) error {
	for _, field := range configMeta.Fields {
		field.RuleMeta = p.parseRule(field.Rule)
	}

	return nil
}

func (p *parser) parseRule(value string) mate.Rule {
	if len(value) == 0 {
		return *mate.NoRule
	}

	for _, pattern := range p.patterns {
		if !pattern.MatchString(value) {
			continue
		}

		ss := pattern.FindAllStringSubmatch(value, -1)[0]

		return mate.Rule{
			Origin: value,
			Key:    ss[1],
			Params: ss[2:],
		}
	}

	return *mate.ErrRule
}
