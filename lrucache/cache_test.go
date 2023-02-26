package lrucache

import (
	"math/rand"
	"runtime/debug"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache_empty(t *testing.T) {
	c := New(0)

	c.Set(1, 2)
	_, ok := c.Get(1)
	require.False(t, ok)
}

func TestCache_update(t *testing.T) {
	c := New(1)

	_, ok := c.Get(1)
	require.False(t, ok)

	c.Set(1, 2)
	v, ok := c.Get(1)
	require.True(t, ok)
	require.Equal(t, 2, v)
}

func TestCache_Get(t *testing.T) {
	c := New(5)

	for i := 0; i < 5; i++ {
		c.Set(i, i)
	}

	c.Get(0)
	c.Get(1)

	c.Set(5, 5)
	c.Set(6, 6)

	var keys, values []int
	c.Range(func(key, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true
	})
	require.Equal(t, []int{4, 0, 1, 5, 6}, keys)
	require.Equal(t, []int{4, 0, 1, 5, 6}, values)
}

func TestCache_Clear(t *testing.T) {
	c := New(5)

	for i := 0; i < 10; i++ {
		c.Set(i, i)
	}

	c.Clear()

	for i := 9; i >= 0; i-- {
		c.Set(i, i)
	}

	var keys, values []int
	c.Range(func(key, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		return true
	})

	require.Equal(t, []int{4, 3, 2, 1, 0}, keys)
	require.Equal(t, []int{4, 3, 2, 1, 0}, values)
}

func TestCache_Range(t *testing.T) {
	c := New(5)

	for i := 0; i < 10; i++ {
		c.Set(i, i)
	}

	var keys, values []int
	c.Range(func(key, value int) bool {
		keys = append(keys, key)
		values = append(values, value)
		return key < 8
	})

	require.Equal(t, []int{5, 6, 7, 8}, keys)
	require.Equal(t, []int{5, 6, 7, 8}, values)
}

func TestCache_eviction(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	for _, tc := range []struct {
		name       string
		cap        int
		numInserts int
		maxKey     int32
	}{
		{
			name:       "empty",
			cap:        0,
			numInserts: 10,
			maxKey:     10,
		},
		{
			name:       "underused",
			cap:        100,
			numInserts: 90,
			maxKey:     1000,
		},
		{
			name:       "simple",
			cap:        100,
			numInserts: 10000,
			maxKey:     1000,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New(tc.cap)

			keyToValue := make(map[int]int)
			for i := 0; i < tc.numInserts; i++ {
				key := int(r.Int31n(tc.maxKey))
				c.Set(key, i)
				keyToValue[key] = i
			}

			var keys, values []int
			c.Range(func(key, value int) bool {
				require.Equal(t, keyToValue[key], value)
				keys = append(keys, key)
				values = append(values, value)
				return true
			})

			expectedLen := tc.cap
			if len(keyToValue) < tc.cap {
				expectedLen = len(keyToValue)
			}
			require.Len(t, values, expectedLen)
			require.True(t, sort.IntsAreSorted(values), "values: %+v", values)

			for _, k := range keys {
				v, ok := c.Get(k)
				require.True(t, ok)
				require.Equal(t, keyToValue[k], v)
			}
		})
	}
}

func TestCache_GCPercent(t *testing.T) {
	var p = debug.SetGCPercent(100)
	require.Equal(t, p, 100)
}

func BenchmarkCache_Set(b *testing.B) {
	for _, tc := range []struct {
		name string
		cap  int
	}{
		{
			name: "small",
			cap:  1,
		},
		{
			name: "medium",
			cap:  1000,
		},
		{
			name: "large",
			cap:  1000000,
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			c := New(tc.cap)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				c.Set(i%tc.cap, i)
			}
		})
	}
}
