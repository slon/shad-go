package reverse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverse(t *testing.T) {
	for i, tc := range []struct {
		input  string
		output string
	}{
		{input: "", output: ""},
		{input: "x", output: "x"},
		{input: "Hello!", output: "!olleH"},
		{input: "Привет", output: "тевирП"},
		{input: "\r\n", output: "\n\r"},
		{input: "\n\n", output: "\n\n"},
		{input: "\t*", output: "*\t"},
		// NB: Диакритика съехала!
		{input: "möp", output: "p̈om"},
		// NB: Иероглиф развалился!,
		{input: "뢴", output: "ᆫᅬᄅ"},
		{input: "Hello, 世界", output: "界世 ,olleH"},
		{input: "ำ", output: "ำ"},
		{input: "ำำ", output: "ำำ"},
		// NB: Эмоджи распался.
		{input: "👩‍❤️‍💋‍👩", output: "👩‍💋‍️❤‍👩"},
		// NB: Эмоджи распался.
		{input: "🏋🏽‍♀️", output: "️♀\u200d🏽🏋"},
		{input: "🙂", output: "🙂"},
		{input: "🙂🙂", output: "🙂🙂"},
		// NB: DE != ED
		{input: "🇩🇪", output: "🇪🇩"},
		// NB: Флаг распался. :)
		{input: "🏳️‍🌈", output: "🌈‍️🏳"},
	} {
		t.Run(fmt.Sprintf("#%v: %v", i, tc.input), func(t *testing.T) {
			require.Equal(t, tc.output, Reverse(tc.input))
		})
	}
}
