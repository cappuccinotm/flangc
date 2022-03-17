package eval

func (s *Scope) cons(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}

	elemExpr, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	listExpr, err := s.Eval(call.Args[1])
	if err != nil {
		return nil, err
	}

	if listExpr == nil {
		return &List{Values: []Expression{elemExpr}}, nil
	}

	list, ok := listExpr.(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: listExpr.Type()}
	}
	return &List{Values: append([]Expression{elemExpr}, list.Values...)}, nil
}

func (s *Scope) tail(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}

	expr, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	list, ok := expr.(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}
	return &List{Values: list.Values[1:]}, nil
}

func (s *Scope) head(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}

	expr, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	list, ok := expr.(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}
	return list.Values[0], nil
}

func (s *Scope) empty(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}

	expr, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	list, ok := expr.(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}
	return &Boolean{Value: len(list.Values) == 0}, nil
}
