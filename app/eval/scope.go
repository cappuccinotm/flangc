package eval

// Scope describes the environment in which an expression is evaluated.
type Scope struct {
	Parent *Scope
	Vars   map[string]Expression
}

// NewScope creates a new scope with the given parent scope.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Vars:   make(map[string]Expression),
	}
}

// Get returns the value of the expression with the given name.
func (s *Scope) Get(name string) Expression {
	if s.Vars != nil {
		if v, ok := s.Vars[name]; ok {
			return v
		}
	}
	if s.Parent != nil {
		return s.Parent.Get(name)
	}
	return nil
}
