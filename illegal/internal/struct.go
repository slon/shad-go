//go:build !change
// +build !change

package internal

import "fmt"

type Struct struct {
	a int
	b string
}

func (s *Struct) String() string {
	return fmt.Sprintf("%d %s", s.a, s.b)
}
