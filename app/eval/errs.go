package eval

import (
	"fmt"
	"errors"
)

// ErrInvalidArguments is returned when the number of arguments passed to a
// function is invalid.
type ErrInvalidArguments struct {
	expected string
	actual   int
}

// Error returns string representation of the error.
func (e ErrInvalidArguments) Error() string {
	return fmt.Sprintf("expected %s arguments, got %d", e.expected, e.actual)
}

// ErrArgumentType is returned when the type of argument passed to a function
// is invalid.
type ErrArgumentType struct {
	expected string
	actual   string
}

// Error returns string representation of the error.
func (e ErrArgumentType) Error() string {
	return fmt.Sprintf("expected argument of type %s, got %s", e.expected, e.actual)
}

// ErrUndefined is returned when the name of a variable is undefined.
type ErrUndefined struct {
	Name string
}

// Error returns string representation of the error.
func (e ErrUndefined) Error() string { return fmt.Sprintf("undefined %s", e.Name) }

// ErrNotFunction is returned when the name of a function is undefined.
type ErrNotFunction struct {
	Name string
}

// Error returns string representation of the error.
func (e ErrNotFunction) Error() string {
	return fmt.Sprintf("%s is not a function", e.Name)
}

var (
	ErrZeroDivision   = errors.New("zero division")
	ErrInvalidContext = errors.New("statement is illegal in this context")
)
