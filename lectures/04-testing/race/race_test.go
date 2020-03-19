package race

import (
	"sync"
	"testing"
)

func TestRace(t *testing.T) {
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(2)

	var i int
	go func() {
		defer wg.Done()
		i = 0
	}()
	go func() {
		defer wg.Done()
		i = 1
	}()
	_ = i
}
