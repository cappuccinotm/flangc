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
		expr, err := s.Eval(call.Args[0])
		if err != nil {
			return nil, err
		}
		if err = s.SetReturn(expr); err != nil {
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
	fmt.Println(expr.FString())
	return Null{}, nil
}

func (s *Scope) lambda(call *Call) (Expression, error) {
	return call, nil
}

func (s *Scope) prog(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}

	exposureListExpr, ok := call.Args[0].(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}

	bodyExpr, ok := call.Args[1].(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[1].Type()}
	}

	scope := NewScope("prog", s, s.PrintNulls)

	for _, expr := range exposureListExpr.Values {
		id, ok := expr.(*Identifier)
		if !ok {
			return nil, ErrArgumentType{expected: "identifier", actual: expr.Type()}
		}

		v, err := s.GetVar(id.Name, true)
		if err != nil {
			return nil, err
		}

		scope.SetVar(id.Name, v)
	}

	for idx, expr := range bodyExpr.Values {
		if _, err := scope.Eval(expr); err != nil {
			return nil, fmt.Errorf("evaluate expression %d: %w", idx, err)
		}
		callExpr, ok := expr.(*Call)
		if !ok {
			continue
		}
		if callExpr.Name == "return" {
			break
		}
	}

	if scope.Return == nil {
		return Null{}, nil
	}

	return scope.Return, nil
}
