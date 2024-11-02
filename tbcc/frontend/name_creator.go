package frontend

import "fmt"

type NameCreator interface {
	VarName() string
	LabelName(prefix string) string
}

type nameCreator struct {
	varCounter    uint
	labelCounters map[string]uint
}

func NewNameCreator() NameCreator {
	return &nameCreator{
		varCounter:    0,
		labelCounters: make(map[string]uint),
	}
}

func (n *nameCreator) VarName() string {
	varName := fmt.Sprintf("tmp.%d", n.varCounter)
	n.varCounter++
	return varName
}

func (n *nameCreator) LabelName(prefix string) string {
	current, ok := n.labelCounters[prefix]
	if !ok {
		current = 0
	}
	ret := fmt.Sprintf("%s.%d", prefix, current)
	current++
	n.labelCounters[prefix] = current
	return ret
}
