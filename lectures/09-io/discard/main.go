package main

import (
	"io"
	"io/ioutil"
	"strings"
)

func main() {
	_, _ = io.Copy(ioutil.Discard, strings.NewReader("nothing of use"))
}
