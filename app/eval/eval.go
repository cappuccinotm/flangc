package eval

import (
	"fmt"
	"strconv"
	"log"
	"errors"
)

// Function describes a defined function.
type Function struct {
	ArgNames []string
	Body     Expression
}

// Scope is an evaluator for expressions.
type Scope struct {
	Parent     *Scope
	Vars       map[string]Expression
	Funcs      map[string]Function
	Context    Expression
	Return     Expression
	PrintNulls bool
}

// NewScope creates a new evaluator.
func NewScope(parent *Scope, printNulls bool) *Scope {
	return &Scope{
		Parent:     parent,
		Vars:       make(map[string]Expression),
		Funcs:      make(map[string]Function),
		PrintNulls: printNulls,
	}
}

// GetVar returns the value of the expression with the given name.
func (s *Scope) GetVar(name string) (Expression, error) {
	if s.Vars != nil {
		if v, ok := s.Vars[name]; ok {
			return v, nil
		}
	}
	return nil, ErrUndefined{Name: name}
}

// SetVar sets the value of the expression with the given name.
func (s *Scope) SetVar(name string, val Expression) {
	if s.Vars == nil {
		s.Vars = make(map[string]Expression)
	}
	s.Vars[name] = val
}

// SetFunc sets the function in the scope.
func (s *Scope) SetFunc(name string, args []string, body Expression) {
	s.Funcs[name] = Function{ArgNames: args, Body: body}
}

// GetFunc returns the function with the given name.
func (s *Scope) GetFunc(name string) (Function, error) {
	if s.Funcs != nil {
		if f, ok := s.Funcs[name]; ok {
			return f, nil
		}
	}
	if s.Parent != nil {
		return s.Parent.GetFunc(name)
	}
	return Function{}, ErrUndefined{Name: name}
}

// SetBreak sets the return value of the evaluator.
func (s *Scope) SetBreak() error {
	ctx, ok := s.Context.(*Call)
	for !ok && s.Parent != nil && ctx.Name != "while" && ctx.Name != "func" {
		s = s.Parent
		ctx, ok = s.Context.(*Call)
	}
	if !ok || s.Parent == nil || ctx.Name != "while" {
		return ErrInvalidContext
	}
	s.Return = brk{}
	return nil
}

// SetReturn sets the return value of the evaluator.
func (s *Scope) SetReturn(val Expression) error {
	ctx, ok := s.Context.(*Call)
	for !ok && s.Parent != nil && ctx.Name != "func" {
		s = s.Parent
		ctx, ok = s.Context.(*Call)
	}
	if !ok || s.Parent == nil {
		return ErrInvalidContext
	}
	s.Return = val
	return nil
}

// Eval evaluates the given expression.
func (s *Scope) Eval(expr Expression) (Expression, error) {
	switch expr := expr.(type) {
	case *Call:
		result, err := s.call(expr)
		if err != nil {
			return nil, fmt.Errorf("call %q: %w", expr.Name, err)
		}
		return result, nil
	case *Number:
		return expr, nil
	case *Identifier:
		v, err := s.GetVar(expr.Name)
		if errors.Is(err, &ErrUndefined{}) {
			fn, err := s.GetFunc(expr.Name)
			if err != nil {
				return nil, err
			}
			return makeLambdaCall(fn), nil
		}
		return v, nil
	case *List:
		return expr, nil
	case *Boolean:
		return expr, nil
	}
	return nil, nil
}

func makeLambdaCall(fn Function) Expression {
	argsList := &List{Values: make([]Expression, len(fn.ArgNames))}

	for _, argname := range fn.ArgNames {
		argsList.Values = append(argsList.Values, &Identifier{Name: argname})
	}

	return &Call{
		Name: "lambda",
		Args: []Expression{argsList, fn.Body},
	}
}

func (s *Scope) call(call *Call) (Expression, error) {
	log.Printf("[DEBUG] call %s", call)
	if expr, ok := builtinMethods[call.Name]; ok {
		return expr(s, call)
	}

	fn, err := s.GetFunc(call.Name)
	if err != nil {
		return nil, err
	}

	if len(call.Args) != len(fn.ArgNames) {
		return nil, ErrInvalidArguments{
			expected: strconv.Itoa(len(fn.ArgNames)),
			actual:   len(call.Args),
		}
	}

	scope := NewScope(s, s.PrintNulls)
	for idx, arg := range fn.ArgNames {
		if nestedCall, ok := call.Args[idx].(*Call); ok && nestedCall.Name == "lambda" {
			argNames, body, err := makeLambdaFunc(nestedCall)
			if err != nil {
				return nil, fmt.Errorf("invalid lambda at argument %d: %w", idx, err)
			}
			scope.SetFunc(arg, argNames, body)
			continue
		}

		argval, err := s.Eval(call.Args[idx])
		if err != nil {
			return nil, fmt.Errorf("evaluate argument %d: %w", idx, err)
		}
		scope.SetVar(arg, argval)
	}

	result, err := scope.Eval(fn.Body)
	if err != nil {
		return nil, fmt.Errorf("evaluate function %s body: %w", call.Name, err)
	}

	return result, nil
}

func makeLambdaFunc(call *Call) ([]string, Expression, error) {
	if len(call.Args) != 2 {
		return nil, nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}

	argListExpr, ok := call.Args[0].(*List)
	if !ok {
		return nil, nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}

	argnames := make([]string, len(argListExpr.Values))
	for idx, expr := range argListExpr.Values {
		identifier, ok := expr.(*Identifier)
		if !ok {
			return nil, nil, ErrArgumentType{expected: "identifier", actual: expr.Type()}
		}

		argnames[idx] = identifier.Name
	}

	return argnames, call.Args[1], nil
}
