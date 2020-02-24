// +build solution

package utf8spacecollapse

import (
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	res := make([]byte, len(input))
	pos := 0
	lastWasSpace := false
	for len(input) > 0 {
		r, n := utf8.DecodeRuneInString(input)
		input = input[n:]
		if unicode.IsSpace(r) {
			if lastWasSpace {
				continue
			}
			res[pos] = ' '
			pos++
			lastWasSpace = true
		} else {
			pos += utf8.EncodeRune(res[pos:], r)
			lastWasSpace = false
		}
	}
	return string(res[:pos])
}
