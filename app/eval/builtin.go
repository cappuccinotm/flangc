package eval

var builtinMethods map[string]func(*Scope, *Call) (Expression, error)

func init() {
	builtinMethods = map[string]func(*Scope, *Call) (Expression, error){
		"quote":    (*Scope).quote,
		"equal":    (*Scope).equal,
		"nonequal": (*Scope).nonequal,
		// list
		"head":  (*Scope).head,
		"tail":  (*Scope).tail,
		"cons":  (*Scope).cons,
		"empty": (*Scope).empty,
		// arithmetic
		"times":     (*Scope).times,
		"plus":      (*Scope).plus,
		"minus":     (*Scope).minus,
		"divide":    (*Scope).div,
		"less":      (*Scope).less,
		"lesseq":    (*Scope).lesseq,
		"greater":   (*Scope).greater,
		"greatereq": (*Scope).greatereq,
		// logic
		"and": (*Scope).and,
		"or":  (*Scope).or,
		"not": (*Scope).not,
		"xor": (*Scope).xor,
		// predecates
		"isnull": is("null"),
		"isbool": is("boolean"),
		"islist": is("list"),
		"isnum":  is("number"),
		// state-related
		"setq": (*Scope).setq,
		"func": (*Scope).setfn,
		// execution flow
		"cond":   (*Scope).cond,
		"while":  (*Scope).while,
		"break":  (*Scope).brk,
		"return": (*Scope).ret,
		"print":  (*Scope).Print,
	}
}

func (s *Scope) quote(call *Call) (Expression, error) {
	return &List{Values: call.Args}, nil
}

func is(typ string) func(*Scope, *Call) (Expression, error) {
	return func(_ *Scope, call *Call) (Expression, error) {
		if len(call.Args) != 1 {
			return nil, ErrInvalidArguments{expected: "1", actual: len(call.Args)}
		}

		if typ == "null" {
			i, ok := call.Args[0].(*Identifier)
			if !ok {
				return &Boolean{Value: false}, nil
			}
			return &Boolean{Value: i.Name == "null"}, nil
		}

		return &Boolean{Value: call.Args[0].Type() == typ}, nil
	}
}

func (s *Scope) nonequal(call *Call) (Expression, error) {
	b, err := s.equal(call)
	if err != nil {
		return b, err
	}
	return &Boolean{Value: !b.(*Boolean).Value}, nil
}

func (s *Scope) equal(call *Call) (Expression, error) {
	if len(call.Args) != 2 {
		return nil, ErrInvalidArguments{expected: "2", actual: len(call.Args)}
	}
	a, err := s.Eval(call.Args[0])
	if err != nil {
		return nil, err
	}
	b, err := s.Eval(call.Args[1])
	if err != nil {
		return nil, err
	}
	return &Boolean{Value: a.Equal(b)}, nil
}
