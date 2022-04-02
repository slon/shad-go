//go:build race

package batcher

func init() {
	race = true
}
