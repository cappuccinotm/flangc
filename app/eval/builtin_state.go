package eval

func setq(eval bool) func(*Scope, *Call) (Expression, error) {
	return func(s *Scope, call *Call) (Expression, error) {
		if len(call.Args) != 2 {
			return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
		}

		name, ok := call.Args[0].(*Identifier)
		if !ok {
			return nil, ErrArgumentType{expected: "identifier", actual: call.Args[0].Type()}
		}

		var (
			val = call.Args[1]
			err error
		)
		if eval {
			if val, err = s.Eval(call.Args[1]); err != nil {
				return nil, err
			}
		}

		s.Set(name.Name, val)

		return nil, nil
	}
}

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

	return nil, nil
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

	return nil, nil
}

func (s *Scope) brk(*Call) (Expression, error) {
	call, ok := s.Context.(*Call)
	if !ok || (call.Name != "while" && call.Name != "prog") {
		return nil, ErrInvalidContext
	}
	if err := s.SetBreak(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Scope) ret(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}
	if err := s.SetReturn(call.Args[1]); err != nil {
		return nil, err
	}
	return nil, nil
}
