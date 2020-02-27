package geometry

import (
	"fmt"
	"testing"
)

func TestPoint(t *testing.T) {
	p := Point{1, 2}
	q := Point{4, 6}
	fmt.Println(Distance(p, q)) // "5", function call
	fmt.Println(p.Distance(q))  // "5", method call
}
