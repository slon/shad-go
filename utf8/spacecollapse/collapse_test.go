package spacecollapse

import (
	"fmt"
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
	} {
		t.Run(fmt.Sprintf("#%v: %v", i, tc.input), func(t *testing.T) {
			require.Equal(t, tc.output, CollapseSpaces(tc.input))
		})
	}
}
