package main

import (
	"fmt"
	"strings"
)

func concat(x, y string) string {
	var builder strings.Builder
	builder.Grow(len(x) + len(y)) // only this line allocates
	builder.WriteString(x)
	builder.WriteString(y)
	return builder.String()
}

func main() {
	fmt.Println(concat("hello ", "world"))
}
