package subpkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// min coverage: 100%

func TestAddOne(t *testing.T) {
	require.Equal(t, 1, AddOne(0))
}
