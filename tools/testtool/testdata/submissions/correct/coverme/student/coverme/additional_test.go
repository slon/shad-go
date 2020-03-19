package coverme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSum2(t *testing.T) {
	require.Equal(t, int64(202), Sum(200, 2))
}
