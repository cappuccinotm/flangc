package eval

func (s *Scope) cons(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}
	list, ok := call.Args[1].(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[1].Type()}
	}
	return &List{Values: append([]Expression{call.Args[0]}, list.Values...)}, nil
}

func (s *Scope) tail(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}
	list, ok := call.Args[0].(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}
	return &List{Values: list.Values[1:]}, nil
}

func (s *Scope) head(call *Call) (Expression, error) {
	if len(call.Args) != 1 {
		return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
	}
	list, ok := call.Args[0].(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[0].Type()}
	}
	return list.Values[0], nil
}
