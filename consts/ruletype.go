package consts

type RuleType string

const (
	NoRule      RuleType = ""
	KeyRule     RuleType = "key"
	UniqueRule  RuleType = "unique"
	RequireRule RuleType = "require"
	RangeRule   RuleType = "range"
	LengthRule  RuleType = "length"
	DecimalRule RuleType = "decimal"
)
