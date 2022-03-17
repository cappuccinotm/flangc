package eval

import (
	"strings"
	"fmt"
	"log"
	"strconv"
)

// Expression describes any language expression.
type Expression interface {
	String() string
	Type() string
	Equal(e Expression) bool
	FString() string
}

// Call represents a function call.
type Call struct {
	Name string
	Args []Expression
}

// FString returns the F language representation of the call.
func (c *Call) FString() string {
	args := make([]string, len(c.Args))
	for idx := range c.Args {
		args[idx] = c.Args[idx].FString()
	}
	return fmt.Sprintf("(%s %s)", c.Name, strings.Join(args, ", "))
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

// Equal returns true if the two calls are equal.
func (*Call) Equal(Expression) bool {
	log.Printf("[WARN] called Call.Equal()")
	return false
}

// Identifier represents an identifier.
type Identifier struct{ Name string }

// FString returns the F language representation of the identifier.
func (i *Identifier) FString() string { panic("must never be called") }

// Type returns the type of the identifier.
func (i *Identifier) Type() string { return "identifier" }

// String returns the string representation of the identifier.
func (i *Identifier) String() string { return i.Name }

// Equal returns true if the two identifiers are equal.
func (*Identifier) Equal(Expression) bool {
	log.Printf("[WARN] called Identifier.Equal()")
	return false
}

// List represents a list of expressions.
type List struct{ Values []Expression }

// Type returns the type of the list.
func (l *List) Type() string { return "list" }

// String returns the string representation of the list.
func (l *List) String() string {
	var args []string
	for _, arg := range l.Values {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(args, ", "))
}

// FString returns the F language representation of the list.
func (l *List) FString() string {
	var args []string
	for _, arg := range l.Values {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("'(%s)", strings.Join(args, " "))
}

// Equal returns true if the two lists are equal.
func (l *List) Equal(b Expression) bool {
	if b, ok := b.(*List); ok {
		if len(l.Values) != len(b.Values) {
			return false
		}
		for i, v := range l.Values {
			if !v.Equal(b.Values[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// Number represents a number.
type Number struct{ Value float64 }

// Type returns the type of the number.
func (n *Number) Type() string { return "number" }

// String returns the string representation of the number.
func (n *Number) String() string { return strconv.FormatFloat(n.Value, 'f', -1, 64) }

// Equal returns true if the two numbers are equal.
func (n *Number) Equal(b Expression) bool {
	if b, ok := b.(*Number); ok {
		return n.Value == b.Value
	}
	return false
}

// FString returns the F language representation of the number.
func (n *Number) FString() string { return n.String() }

// Boolean represents a boolean.
type Boolean struct{ Value bool }

// Type returns the type of the boolean.
func (b *Boolean) Type() string { return "boolean" }

// String returns the string representation of the boolean.
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

// Equal returns true if the two booleans are equal.
func (b *Boolean) Equal(b2 Expression) bool {
	if b2, ok := b2.(*Boolean); ok {
		return b.Value == b2.Value
	}
	return false
}

// FString returns the F language representation of the boolean.
func (b *Boolean) FString() string { return b.String() }

// Null represents a null value.
type Null struct{}

func (n Null) String() string { return "null" }

// Type returns the type of the null.
func (n Null) Type() string { return "null" }

// Equal returns true if the two nulls are equal.
func (n Null) Equal(e Expression) bool {
	if _, ok := e.(Null); ok {
		return true
	}
	return false
}

// FString returns the F language representation of the null.
func (n Null) FString() string { return "null" }

type brk struct{}

func (b brk) FString() string         { panic("must never be called") }
func (b brk) String() string          { panic("must never be called") }
func (b brk) Type() string            { panic("must never be called") }
func (b brk) Equal(e Expression) bool { panic("must never be called") }
