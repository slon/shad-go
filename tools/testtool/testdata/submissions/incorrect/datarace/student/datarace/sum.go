//go:build !solution
// +build !solution

package datarace

import "sync"

func Sum(a, b int64) int64 {
	var s int64

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s += a
	}()

	go func() {
		defer wg.Done()
		s += b
	}()

	wg.Wait()
	return s
}
