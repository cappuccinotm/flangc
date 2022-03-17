package eval

import "fmt"

func (s *Scope) setq(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}

	name, ok := call.Args[0].(*Identifier)
	if !ok {
		return nil, ErrArgumentType{expected: "identifier", actual: call.Args[0].Type()}
	}

	val, err := s.Eval(call.Args[1])
	if err != nil {
		return nil, err
	}

	s.SetVar(name.Name, val)

	return nil, nil
}

func (s *Scope) setfn(call *Call) (Expression, error) {
	if len(call.Args) != 3 {
		return nil, ErrInvalidArguments{expected: "3", actual: len(call.Args)}
	}

	name, ok := call.Args[0].(*Identifier)
	if !ok {
		return nil, ErrArgumentType{expected: "identifier", actual: call.Args[0].Type()}
	}

	argList, ok := call.Args[1].(*List)
	if !ok {
		return nil, ErrArgumentType{expected: "list", actual: call.Args[1].Type()}
	}

	args := make([]string, len(argList.Values))
	for idx, expr := range argList.Values {
		str, ok := expr.(*Identifier)
		if !ok {
			return nil, fmt.Errorf("argument %d: %w", idx, ErrArgumentType{expected: "identifier", actual: expr.Type()})
		}
		args[idx] = str.Name
	}

	s.SetFunc(name.Name, args, call.Args[2])
	return nil, nil
}
