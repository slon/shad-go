// +build !solution

package lrucache

import (
	"container/list"
)

type Var struct {
	key		int
	value	int
}

type LRUCache struct {
	data     map[int]*list.Element
	queue    *list.List
	capacity int
	size     int
}

func (cache *LRUCache) Set(key, value int) {
	if cache.capacity == 0 {
		return
	}

	if v, ok := cache.data[key]; !ok {

		if cache.capacity == cache.size {
			oldest := cache.queue.Back().Value.(*Var)
			delete(cache.data, oldest.key)

			cache.queue.Remove(cache.queue.Back())
			cache.queue.PushFront(&Var{key, value})
			cache.data[key] = cache.queue.Front()
		} else {
			cache.queue.PushFront(&Var{key, value})
			cache.data[key] = cache.queue.Front()
			cache.size++
		}
	} else {
		cache.queue.MoveToFront(v)
		cache.queue.Front().Value.(*Var).value = value
	}
}

func (cache *LRUCache) Get(key int) (value int, has bool) {
	val, has := cache.data[key]
	if !has {
		return
	}

	cache.queue.MoveToFront(val)
	return val.Value.(*Var).value, has
}

func (cache *LRUCache) Clear() {
	cache.size = 0
	cache.queue = list.New()
	cache.data = make(map[int]*list.Element, cache.capacity)
}

func (cache *LRUCache) Range(f func(key, value int) bool) {
	for e := cache.queue.Back(); e != nil; e = e.Prev() {
		elem := e.Value.(*Var)
		if !f(elem.key, elem.value) {
			return
		}
	}
}

func (cache *LRUCache) Init(cap int) *LRUCache {
	cache.data = make(map[int]*list.Element, cache.capacity)
	cache.queue = list.New()
	cache.capacity = cap
	cache.size = 0
	return cache
}

func New(cap int) Cache {
	return new(LRUCache).Init(cap)
}
