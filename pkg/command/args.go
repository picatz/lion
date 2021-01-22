package command

import (
	"fmt"
	"strings"
)

type Arg struct {
	Index int
	Name  string
	Desc  string
	Value interface{}
}

func (c Arg) StringValue() (string, error) {
	v, ok := c.Value.(string)
	if !ok {
		return "", fmt.Errorf("argument %q at index %d is %T not a string", c.Name, c.Index, v)
	}
	return v, nil
}

type Args map[string]*Arg

func (c Args) HelpString() string {
	all := []string{}
	for k := range c {
		all = append(all, fmt.Sprintf("<%s>", k))
	}
	return strings.Join(all, " ")
}

func (c Args) ForIndex(index int) (*Arg, error) {
	for _, a := range c {
		if a.Index == index {
			return a, nil
		}
	}
	return nil, fmt.Errorf("no command argument for index %d", index)
}

func (c Args) StringValueForIndex(index int) (string, error) {
	for _, a := range c {
		if a.Index == index {
			return a.StringValue()
		}
	}
	return "", fmt.Errorf("no command argument for index %d", index)
}
