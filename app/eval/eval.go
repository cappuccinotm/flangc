package eval

// Scope is an evaluator for expressions.
type Scope struct {
	Parent  *Scope
	Vars    map[string]Expression
	Context Expression
	Return  Expression
}

// NewScope creates a new evaluator.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Vars:   make(map[string]Expression),
	}
}

// Get returns the value of the expression with the given name.
func (s *Scope) Get(name string) (Expression, error) {
	if s.Vars != nil {
		if v, ok := s.Vars[name]; ok {
			return v, nil
		}
	}
	return nil, ErrUndefined{Name: name}
}

// Set sets the value of the expression with the given name.
func (s *Scope) Set(name string, val Expression) {
	if s.Vars == nil {
		s.Vars = make(map[string]Expression)
	}
	s.Vars[name] = val
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

	return nil, nil
}
