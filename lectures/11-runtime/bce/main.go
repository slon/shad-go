package main

import "encoding/binary"

func main() {
	d := make([]byte, 16)

	// if 10 < len(d) {
	// 	panic("out of bound")
	// }
	d[10] = 12

	i := binary.LittleEndian.Uint64(d)
	_ = i
}
