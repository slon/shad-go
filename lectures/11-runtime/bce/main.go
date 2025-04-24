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

func Uint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}
