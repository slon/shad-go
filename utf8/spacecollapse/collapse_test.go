package spacecollapse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollapseSpaces(t *testing.T) {
	for i, tc := range []struct {
		input  string
		output string
	}{
		{input: "", output: ""},
		{input: "x", output: "x"},
		{input: "Hello,   World!", output: "Hello, World!"},
		{input: "Привет,\tМир!", output: "Привет, Мир!"},
		{input: "\r\n", output: " "},
		{input: "\n\n", output: " "},
		{input: "\t*", output: " *"},
		{input: " \t \t ", output: " "},
		{input: " \tx\t ", output: " x "},
		{input: "\xff\x00   \xff\x00", output: "\xef\xbf\xbd\x00 \xef\xbf\xbd\x00"},
	} {
		t.Run(fmt.Sprintf("#%v: %v", i, tc.input), func(t *testing.T) {
			require.Equal(t, tc.output, CollapseSpaces(tc.input))
		})
	}
}

func BenchmarkCollapse(b *testing.B) {
	input := strings.Repeat("🙂  🙂", 100)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = CollapseSpaces(input)
	}
}
