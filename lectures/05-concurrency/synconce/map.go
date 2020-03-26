package synconce

import "sync"

var cache sync.Map

type result struct{}

func do(key string) *result { return new(result) }

type entry struct {
	res *result
	sync.Once
}

func get(key string) *result {
	myEntry := &entry{}

	old, loaded := cache.LoadOrStore(key, myEntry)
	if loaded {
		myEntry = old.(*entry)
	}

	myEntry.Do(func() {
		myEntry.res = do(key)
	})

	return myEntry.res
}

// OMIT
