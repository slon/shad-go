package genericsum_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.com/slon/shad-go/genericsum"
)

func TestMin(t *testing.T) {
	assert.Equal(t, 10, genericsum.Min(10, 20))
	assert.Equal(t, -10, genericsum.Min(20, -10))

	assert.Equal(t, 1.0, genericsum.Min(1.0, 2.0))
	assert.Equal(t, int64(10), genericsum.Min(int64(10), 20))

	assert.Equal(t, "abc", genericsum.Min("def", "abc"))
	type myString string
	assert.Equal(t, myString("abc"), genericsum.Min(myString("def"), myString("abc")))
}

func TestSortSlice(t *testing.T) {
	t.Run("ints", func(t *testing.T) {
		inputs, expected := []int{3, 6, 2, 4, 5, 1}, []int{1, 2, 3, 4, 5, 6}
		genericsum.SortSlice(inputs)
		assert.Equal(t, expected, inputs)
	})
	t.Run("ints64", func(t *testing.T) {
		inputs, expected := []int64{3, 6, 2, 4, 5, 1}, []int64{1, 2, 3, 4, 5, 6}
		genericsum.SortSlice(inputs)
		assert.Equal(t, expected, inputs)
	})
	t.Run("strings", func(t *testing.T) {
		inputs, expected := []string{"d", "b", "ab", "a"}, []string{"a", "ab", "b", "d"}
		genericsum.SortSlice(inputs)
		assert.Equal(t, expected, inputs)
	})

	type myStringSlice []string
	t.Run("strings custom type", func(t *testing.T) {
		inputs, expected := myStringSlice([]string{"d", "b", "ab", "a"}), myStringSlice([]string{"a", "ab", "b", "d"})
		genericsum.SortSlice(inputs)
		assert.Equal(t, expected, inputs)
	})
}

func TestMapsEqual(t *testing.T) {
	assert.True(t, genericsum.MapsEqual(map[string]string{"1": "3", "2": "4"}, map[string]string{"2": "4", "1": "3"}))

	var i int
	assert.False(t, genericsum.MapsEqual(map[string]*int{"1": &i, "2": nil}, map[string]*int{"1": &i}))
	assert.False(t, genericsum.MapsEqual(map[string]*int{"1": &i}, map[string]*int{"1": &i, "2": nil}))

	assert.False(t, genericsum.MapsEqual(map[string]*int{"1": new(int)}, map[string]*int{"1": new(int)}),
		"different pointers")

	type k struct {
		i int
		s string
	}
	assert.True(t, genericsum.MapsEqual(map[k]k{{10, "abc"}: {20, "def"}}, map[k]k{{10, "abc"}: {20, "def"}}))

	type myMap map[int]int
	assert.True(t, genericsum.MapsEqual(myMap(nil), myMap(nil)), "type aliases must also be ok")
}

func TestSliceContains(t *testing.T) {
	assert.True(t, genericsum.SliceContains([]int{5, 9, 12}, 5))
	assert.False(t, genericsum.SliceContains([]int{5, 9, 12}, 7))

	type k struct{ i, j int }
	assert.True(t, genericsum.SliceContains([]k{{1, 2}, {5, 7}, {9, 12}}, k{5, 7}))

	type mySlice []k
	assert.True(t, genericsum.SliceContains(mySlice{{1, 2}, {5, 7}, {9, 12}}, k{5, 7}))
}

func TestMergeChansTypes(t *testing.T) {
	t.Run("floats", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		chans := make([]chan float64, 5)
		chanArgs := make([]<-chan float64, 5)
		for i := range chans {
			chans[i] = make(chan float64, 1)
			chanArgs[i] = chans[i]
		}

		ch := genericsum.MergeChans(chanArgs...)

		for _, ch := range chans {
			close(ch)
		}

		for range ch {
		}
	})

	t.Run("structs", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		type tp struct{}

		chans := make([]chan tp, 5)
		chanArgs := make([]<-chan tp, 5)
		for i := range chans {
			chans[i] = make(chan tp, 1)
			chanArgs[i] = chans[i]
		}

		ch := genericsum.MergeChans(chanArgs...)

		for _, ch := range chans {
			close(ch)
		}

		for range ch {
		}
	})
}

func TestMergeChans(t *testing.T) {
	defer goleak.VerifyNone(t)

	const numChans = 10

	chans := make([]chan int, numChans)
	chanArgs := make([]<-chan int, numChans)
	for i := range chans {
		chans[i] = make(chan int, 1)
		chanArgs[i] = chans[i]
	}

	ch := genericsum.MergeChans(chanArgs...)
	chans[5] <- 10
	assert.Equal(t, 10, <-ch)
	chans[1] <- 5
	assert.Equal(t, 5, <-ch)

	const numIter = 1000
	receivedNumbers := make([]bool, numIter)

	go func() {
		for i := 0; i < numIter; i++ {
			chans[i%numChans] <- i
		}
		for _, ch := range chans {
			close(ch)
		}
	}()

	// don't ever deadlock
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var receivedNum int

Loop:
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				break Loop
			}
			require.True(t, v >= 0 && v < len(receivedNumbers), "invalid received number", v)
			require.False(t, receivedNumbers[v], "shouldn't receive number twice", v)
			receivedNumbers[v] = true
			receivedNum++
		case <-ctx.Done():
			require.FailNow(t, "timeouted")
		}
	}
	assert.Equal(t, numIter, receivedNum)
}

func TestIsHermitianMatrix(t *testing.T) {
	assert.True(t, genericsum.IsHermitianMatrix([][]int{
		{1, 7, 9},
		{7, 2, 12},
		{9, 12, 19},
	}))
	assert.False(t, genericsum.IsHermitianMatrix([][]int{
		{1, 12, 8},
		{3, 4, 7},
		{8, 7, 11},
	}))
	assert.True(t, genericsum.IsHermitianMatrix([][]float32{
		{1.0, 7.0, 9.0},
		{7.0, 2.0, 12.0},
		{9.0, 12.0, 19.0},
	}))
	assert.False(t, genericsum.IsHermitianMatrix([][]float32{
		{1.0, 12.0, 8.0},
		{3.0, 4.0, 7.0},
		{8.0, 7.0, 11.0},
	}))
	assert.True(t, genericsum.IsHermitianMatrix([][]complex64{
		{1, 3 + 2i},
		{3 - 2i, 4},
	}))
	assert.True(t, genericsum.IsHermitianMatrix([][]complex128{
		{1, 3 + 2i, 9 - 1i},
		{3 - 2i, 5, 7 - 3i},
		{9 + 1i, 7 + 3i, 19},
	}))
	assert.False(t, genericsum.IsHermitianMatrix([][]complex64{
		{1 + 1i, 3 + 2i},
		{3 - 2i, 4},
	}))

	first, second := [][]int{{1, 2}, {3, 4}}, [][]int{{1, 2}, {3, 4}}
	assert.False(t, genericsum.IsHermitianMatrix(first))
	assert.Equal(t, second, first, "shouldn't change matrix in the method")
}
