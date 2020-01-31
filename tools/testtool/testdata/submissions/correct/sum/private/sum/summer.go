// +build !change

package sum

type Summer interface {
	Sum(a, b int64) int64
}

// Summer implementation.
type summer struct{}

func (s *summer) Sum(a, b int64) int64 {
	return Sum(a, b)
}
