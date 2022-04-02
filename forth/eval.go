//go:build !solution

package main

type Evaluator struct {
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	return nil, nil
}
