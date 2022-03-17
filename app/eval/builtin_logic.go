package eval

func (s *Scope) and(call *Call) (Expression, error) {
	arg1, arg2, err := s.castBooleanArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Boolean{Value: arg1.Value && arg2.Value}, nil
}

func (s *Scope) or(call *Call) (Expression, error) {
	arg1, arg2, err := s.castBooleanArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Boolean{Value: arg1.Value || arg2.Value}, nil
}

func (s *Scope) xor(call *Call) (Expression, error) {
	arg1, arg2, err := s.castBooleanArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Boolean{Value: (arg1.Value || arg2.Value) && !(arg1.Value && arg2.Value)}, nil
}

func (s *Scope) not(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}
	expr, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}

	arg, ok := expr.(*Boolean)
	if !ok {
		return nil, ErrArgumentType{expected: "boolean", actual: expr.Type()}
	}

	return &Boolean{Value: !arg.Value}, nil
}

func (s *Scope) castBooleanArguments(exprs []Expression) (*Boolean, *Boolean, error) {
	if len(exprs) != 2 {
		return nil, nil, ErrInvalidArguments{expected: "2", actual: len(exprs)}
	}
	a, err := s.Eval(exprs[0])
	if err != nil {
		return nil, nil, err
	}
	b, err := s.Eval(exprs[1])
	if err != nil {
		return nil, nil, err
	}
	arg1, ok := a.(*Boolean)
	if !ok {
		return nil, nil, ErrArgumentType{expected: "number", actual: a.Type()}
	}
	arg2, ok := b.(*Boolean)
	if !ok {
		return nil, nil, ErrArgumentType{expected: "number", actual: b.Type()}
	}
	return arg1, arg2, nil
}
