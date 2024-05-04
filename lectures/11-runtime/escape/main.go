package main

var g any

func main() {
	p0 := make([]byte, 10)
	p0[0] = 'f'

	var p1 [10]byte

	copy(p1[:], p0)

	g = p1
}
