package main

import "fmt"

func F[T ~[]T](t T) T {
	return t[1][3][3][7][6][6][6]
}

type G []G

func main() {
	g := make(G, 10)
	for i := range g {
		g[i] = g
	}
	fmt.Println(F(g))
}
