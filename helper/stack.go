package helper

type Stack struct {
	elements []any
	size     int
}

func (p *Stack) Push(v any) {
	p.elements = append(p.elements, v)
	p.size++
}

func (p *Stack) Pop() (any, bool) {
	if p.IsEmpty() {
		return nil, false
	}

	last := p.size - 1
	element := p.elements[last]

	p.elements[last] = nil
	p.elements = p.elements[:last]
	p.size--

	return element, true
}

func (p *Stack) IsEmpty() bool {
	return p.size == 0
}
