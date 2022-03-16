package parser

import (
	"strings"
	"fmt"
)

type Expression interface {
	String() string
}

type Call struct {
	Name string
	Args []Expression
}

func (c *Call) String() string {
	var args []string
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(args, ", "))
}

type Identifier struct {
	Name string
}

func (i *Identifier) String() string { return i.Name }

type List struct {
	Elements []Expression
}

func (l *List) String() string {
	return fmt.Sprintf("%v", l.Elements)
}

type Number struct {
	Value float64
}

func (n *Number) String() string {
	return fmt.Sprintf("%f", n.Value)
}
