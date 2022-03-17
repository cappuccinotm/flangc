package eval

import "fmt"

func (s *Scope) cond(call *Call) (Expression, error) {
	if len(call.Args) < 2 || len(call.Args) > 3 {
		return nil, ErrInvalidArguments{expected: "2 or 3", actual: len(call.Args)}
	}

	predicate, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	b, ok := predicate.(*Boolean)
	if !ok {
		return nil, ErrArgumentType{expected: "boolean", actual: predicate.Type()}
	}

	if b.Value {
		return s.Eval(call.Args[1])
	}

	if len(call.Args) == 3 {
		return s.Eval(call.Args[2])
	}

	return Null{}, nil
}

func (s *Scope) while(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}

	predicate, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	b, ok := predicate.(*Boolean)
	if !ok {
		return nil, ErrArgumentType{expected: "boolean", actual: predicate.Type()}
	}

	for b.Value && s.Return == nil {
		return s.Eval(call.Args[1])
	}

	return Null{}, nil
}

func (s *Scope) brk(*Call) (Expression, error) {
	call, ok := s.Context.(*Call)
	if !ok || (call.Name != "while" && call.Name != "prog") {
		return nil, ErrInvalidContext
	}
	if err := s.SetBreak(); err != nil {
		return nil, err
	}
	return Null{}, nil
}

func (s *Scope) ret(call *Call) (Expression, error) {
	if len(call.Args) > 1 {
		return nil, ErrInvalidArguments{expected: "0 or 1", actual: len(call.Args)}
	}
	if len(call.Args) == 1 {
		if err := s.SetReturn(call.Args[1]); err != nil {
			return nil, err
		}
	}

	return Null{}, nil
}

func (s *Scope) Print(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}

	expr, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	if expr == nil || expr.String() == "null" {
		if s.PrintNulls {
			fmt.Println("null")
		}
		return Null{}, nil
	}
	fmt.Println(expr.String())
	return Null{}, nil
}

func (s *Scope) lambda(call *Call) (Expression, error) {
	return call, nil
}
