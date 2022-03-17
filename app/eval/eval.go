package eval

// Evaluator is an evaluator for expressions.
type Evaluator struct {
	scope *Scope
}

// NewEvaluator creates a new evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{scope: NewScope(nil)}
}

func (e *Evaluator) Eval(expr Expression) (Expression, error) {
	return nil, nil
}
