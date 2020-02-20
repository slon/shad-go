package main

import "fmt"

func f() int      { return 0 }
func g(x int) int { return x }

func example() {
	if x := f(); x == 0 {
		fmt.Println(x)
	} else if y := g(x); x == y {
		fmt.Println(x, y)
	} else {
		fmt.Println(x, y)
	}
}
