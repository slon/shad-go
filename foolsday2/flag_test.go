// +build private

package main

import (
	"fmt"
	"os"
	"testing"
)

const flag = "FLAG{foolsday2:hidden-gem:07a5e6469f2178616cba4e9a0410e050}"

func TestHiddenGem(t *testing.T) {
	fmt.Fprintf(os.Stderr, "Here's your flag: %s", flag)
}
