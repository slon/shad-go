package allocs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounter_Count(t *testing.T) {
	repeats := 10
	data := strings.Repeat("a b c\n", repeats)
	data = data[:len(data)-1]
	r := strings.NewReader(data)

	c := NewEnhancedCounter()

	err := c.Count(r)
	require.NoError(t, err)

	expected := fmt.Sprintf("word 'a' has %d occurrences\n", repeats)
	expected += fmt.Sprintf("word 'b' has %d occurrences\n", repeats)
	expected += fmt.Sprintf("word 'c' has %d occurrences\n", repeats)
	actual := c.String()
	require.Equal(t, expected, actual)
}

func Benchmark(b *testing.B) {
	repeats := 10000
	data := strings.Repeat("a b c d e f g h i j k l m n o p q r s t u v w x y z\n", repeats)
	data = data[:len(data)-1]

	b.Run("count", func(b *testing.B) {
		r := strings.NewReader(data)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c := NewEnhancedCounter()
			_ = c.Count(r)
		}
	})

	b.Run("main", func(b *testing.B) {
		r := strings.NewReader(data)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c := NewEnhancedCounter()
			_ = c.Count(r)
			_ = c.String()
		}
	})
}
