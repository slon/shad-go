//go:build race
// +build race

package batcher

func init() {
	race = true
}
