package mate

var (
	NoRule  = &Rule{}
	ErrRule = &Rule{}
)

type Rule struct {
	Origin string
	Key    string
	Params []string
}
