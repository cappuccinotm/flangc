package eval

import (
	"strings"
	"fmt"
	"log"
)

// Expression describes any language expression.
type Expression interface {
	String() string
	Type() string
	//Equal(other Expression) bool
}

// Call represents a function call.
type Call struct {
	Name string
	Args []Expression
}

// Type returns the type of the call.
func (c *Call) Type() string { return "function call" }

// String returns the string representation of the call.
func (c *Call) String() string {
	var args []string
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(args, ", "))
}

// Equal returns true if the call is equal to the other call.
func (*Call) Equal(Expression) bool {
	log.Printf("[WARN] equal on function calls is not supported")
	return false
}

// Identifier represents an identifier.
type Identifier struct{ Name string }

// Type returns the type of the identifier.
func (i *Identifier) Type() string { return "identifier" }

// String returns the string representation of the identifier.
func (i *Identifier) String() string { return i.Name }

// Equal returns true if the identifier is equal to the other identifier.
func (i *Identifier) Equal(other Expression) bool {
	panic("not implemented")
	return false
}

// List represents a list of expressions.
type List struct{ Values []Expression }

// Type returns the type of the list.
func (l *List) Type() string { return "list" }

// String returns the string representation of the list.
func (l *List) String() string { return fmt.Sprintf("%v", l.Values) }

// Number represents a number.
type Number struct{ Value float64 }

// Type returns the type of the number.
func (n *Number) Type() string { return "number" }

// String returns the string representation of the number.
func (n *Number) String() string { return fmt.Sprintf("%f", n.Value) }

// Boolean represents a boolean.
type Boolean struct{ Value bool }

// Type returns the type of the boolean.
func (b *Boolean) Type() string { return "boolean" }

// String returns the string representation of the boolean.
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }
