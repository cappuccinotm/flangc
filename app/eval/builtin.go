package eval

import "fmt"

var builtinMethods = map[string]func(*Evaluator, *Call) (Expression, error){
	"times":  (*Evaluator).times,
	"plus":   (*Evaluator).plus,
	"minus":  (*Evaluator).minus,
	"divide": (*Evaluator).div,
	"head":   (*Evaluator).head,
	"tail":   (*Evaluator).tail,
	"cons":   (*Evaluator).cons,
	"equal":  (*Evaluator).equal,
}

func (e *Evaluator) nonequal(call *Call) (Expression, error) {
	b, err := e.equal(call)
	if err != nil {
		return b, err
	}
	return &Boolean{Value: !b.(*Boolean).Value}, nil
}

func (e *Evaluator) equal(call *Call) (Expression, error) {
	//if len(call.Args) != 2 {
	//	return nil, fmt.Errorf("equal takes exactly two arguments")
	//}
	//return &Boolean{Value: call.Args[0].Equal(call.Args[1])}, nil
	panic("not implemented")
}

func (e *Evaluator) cons(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, fmt.Errorf("cons takes two arguments")
	}
	list, ok := call.Args[1].(*List)
	if !ok {
		return nil, fmt.Errorf("cons takes a list as second argument")
	}
	return &List{Values: append([]Expression{call.Args[0]}, list.Values...)}, nil
}

func (e *Evaluator) tail(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, fmt.Errorf("tail: expected 1 argument, got %d", len(call.Args))
	}
	if list, ok := call.Args[0].(*List); ok {
		if len(list.Values) == 0 {
			return nil, fmt.Errorf("tail: empty list")
		}
		return &List{Values: list.Values[1:]}, nil
	}
	return nil, fmt.Errorf("tail: expected list, got %s", call.Args[0].Type())
}

func (e *Evaluator) head(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, fmt.Errorf("head: expected 1 argument, got %d", len(call.Args))
	}
	if list, ok := call.Args[0].(*List); ok {
		if len(list.Values) == 0 {
			return nil, fmt.Errorf("head: empty list")
		}
		return list.Values[0], nil
	}
	return nil, fmt.Errorf("head: expected list, got %s", call.Args[0].Type())
}

func (e *Evaluator) minus(call *Call) (Expression, error) {
	arg1, arg2, err := e.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Number{Value: arg1.Value - arg2.Value}, nil
}

func (e *Evaluator) div(call *Call) (Expression, error) {
	arg1, arg2, err := e.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	if arg2.Value == 0 {
		return nil, errZeroDivision
	}
	return &Number{Value: arg1.Value / arg2.Value}, nil
}

func (e *Evaluator) plus(call *Call) (Expression, error) {
	arg1, arg2, err := e.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Number{Value: arg1.Value + arg2.Value}, nil
}

func (e *Evaluator) times(call *Call) (Expression, error) {
	arg1, arg2, err := e.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Number{Value: arg1.Value * arg2.Value}, nil
}

func (e *Evaluator) castArithmeticArguments(exprs []Expression) (*Number, *Number, error) {
	if len(exprs) != 2 {
		return nil, nil, ErrInvalidArguments{
			expected: 2,
			actual:   len(exprs),
		}
	}
	a, err := e.Eval(exprs[0])
	if err != nil {
		return nil, nil, err
	}
	b, err := e.Eval(exprs[1])
	if err != nil {
		return nil, nil, err
	}
	arg1, ok := a.(*Number)
	if !ok {
		return nil, nil, ErrArgumentType{expected: "number", actual: a.Type()}
	}
	arg2, ok := b.(*Number)
	if !ok {
		return nil, nil, ErrArgumentType{expected: "number", actual: b.Type()}
	}
	return arg1, arg2, nil
}

// ErrInvalidArguments is returned when the number of arguments passed to a
// function is invalid.
type ErrInvalidArguments struct {
	expected int
	actual   int
}

// Error returns string representation of the error.
func (e ErrInvalidArguments) Error() string {
	return fmt.Sprintf("expected %d arguments, got %d", e.expected, e.actual)
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

var errZeroDivision = fmt.Errorf("zero division")
