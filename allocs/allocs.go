//go:build !solution
// +build !solution

package allocs

// implement your Counter below

func NewEnhancedCounter() Counter {
	return NewBaselineCounter()
}
