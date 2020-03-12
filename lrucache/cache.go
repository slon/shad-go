// +build !change

package lrucache

type Cache interface {
	// Get returns value associated with the key.
	//
	// The second value is a bool that is true if the key exists in the cache,
	// and false if not.
	Get(key int) (int, bool)
	// Set updates value associated with the key.
	//
	// If there is no key in the cache new (key, value) pair is created.
	Set(key, value int)
	// Range calls function f on all elements of the cache
	// in increasing access time order.
	//
	// Stops earlier if f returns false.
	Range(f func(key, value int) bool)
	// Clear removes all keys and values from the cache.
	Clear()
}
