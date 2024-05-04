package main

import "unsafe"

const x = unsafe.Sizeof(int(0))

type y [x]int

type s struct {
	a bool
	b int16
	c []int
}

const z = unsafe.Offsetof(s{}.b)

func main() {
	_ = x

	var _ y
}
