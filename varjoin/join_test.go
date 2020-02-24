package varjoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	for _, tc := range []struct {
		sep    string
		args   []string
		result string
	}{
		{
			sep:    "/",
			args:   []string{},
			result: "",
		},
		{
			sep:    "/",
			args:   []string{"1", "2", "3"},
			result: "1/2/3",
		},
		{
			sep:    "",
			args:   []string{"h", "e", "ll", "o"},
			result: "hello",
		},
	} {
		t.Run(tc.result, func(t *testing.T) {
			require.Equal(t, tc.result, Join(tc.sep, tc.args...))
		})
	}
}
