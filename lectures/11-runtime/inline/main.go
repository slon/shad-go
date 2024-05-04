package main

import (
	"os"
)

func F(a int) int {
	return 2 / (a - 1)
}

func main() {
	os.Exit(F(len(os.Args)))
}
