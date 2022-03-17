package eval

func (s *Scope) greatereq(call *Call) (Expression, error) {
	greaterv, err := s.greater(call)
	if err != nil {
		return nil, err
	}

	eqv, err := s.equal(call)
	if err != nil {
		return nil, err
	}

	return &Boolean{Value: greaterv.(*Boolean).Value || eqv.(*Boolean).Value}, nil
}

func (s *Scope) greater(call *Call) (Expression, error) {
	arg1, arg2, err := s.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Boolean{Value: arg1.Value > arg2.Value}, nil
}

func (s *Scope) lesseq(call *Call) (Expression, error) {
	lessv, err := s.less(call)
	if err != nil {
		return nil, err
	}

	eqv, err := s.equal(call)
	if err != nil {
		return nil, err
	}

	return &Boolean{Value: lessv.(*Boolean).Value || eqv.(*Boolean).Value}, nil
}

func (s *Scope) less(call *Call) (Expression, error) {
	arg1, arg2, err := s.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Boolean{Value: arg1.Value < arg2.Value}, nil
}

func (s *Scope) minus(call *Call) (Expression, error) {
	arg1, arg2, err := s.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Number{Value: arg1.Value - arg2.Value}, nil
}

func (s *Scope) div(call *Call) (Expression, error) {
	arg1, arg2, err := s.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	if arg2.Value == 0 {
		return nil, ErrZeroDivision
	}
	return &Number{Value: arg1.Value / arg2.Value}, nil
}

func (s *Scope) plus(call *Call) (Expression, error) {
	arg1, arg2, err := s.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Number{Value: arg1.Value + arg2.Value}, nil
}

func (s *Scope) times(call *Call) (Expression, error) {
	arg1, arg2, err := s.castArithmeticArguments(call.Args)
	if err != nil {
		return nil, err
	}
	return &Number{Value: arg1.Value * arg2.Value}, nil
}

func (s *Scope) castArithmeticArguments(exprs []Expression) (*Number, *Number, error) {
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
