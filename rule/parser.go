package rule

import (
	"regexp"
	"strings"

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
		regexp.MustCompile(`(enum):(\S+(?:,\S+)*)$`),
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

		var (
			ss     = pattern.FindAllStringSubmatch(value, -1)[0]
			key    = ss[1]
			params = ss[2:]
		)

		if key == "enum" {
			params = strings.Split(ss[2], ",")
		}

		return mate.Rule{
			Origin: value,
			Key:    key,
			Params: params,
		}
	}

	return *mate.ErrRule
}
