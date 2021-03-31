// +build private

package main

import "testing"

const flag = "FLAG{foolsday2:hidden-gem:07a5e6469f2178616cba4e9a0410e050}"

func TestHiddenGem(t *testing.T) {
	t.Logf("Here's your flag: %s", flag)
}
