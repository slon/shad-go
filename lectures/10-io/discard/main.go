package main

import (
	"io"
	"strings"
)

func main() {
	_, _ = io.Copy(io.Discard, strings.NewReader("nothing of use"))
}
