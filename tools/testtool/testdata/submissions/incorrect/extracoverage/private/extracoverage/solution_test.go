// +build solution

package extracoverage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	require.Equal(t, int64(2), Sum(1, 1))
	require.Equal(t, int64(1), Sum(0, 1))
}
