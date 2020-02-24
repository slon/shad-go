// +build solution

package utf8reverse

import (
	"unicode/utf8"
)

func Reverse(input string) string {
	rs := []rune{}
	sz := 0
	for len(input) > 0 {
		r, n := utf8.DecodeRuneInString(input)
		rs = append(rs, r)
		sz += utf8.RuneLen(r)
		input = input[n:]
	}
	bs := make([]byte, sz)
	for i, j := 0, 0; i < len(rs); i++ {
		n := utf8.EncodeRune(bs[j:], rs[len(rs)-i-1])
		j += n
	}
	return string(bs)
}
