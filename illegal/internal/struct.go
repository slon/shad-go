//go:build !change

package internal

import "fmt"

type privateType struct {
	x int
}

func NewPrivateType(x int) any {
	return privateType{x}
}

type Struct struct {
	a int
	b string
	p privateType
}

func (s *Struct) String() string {
	return fmt.Sprintf("%d %s %d", s.a, s.b, s.p.x)
}
