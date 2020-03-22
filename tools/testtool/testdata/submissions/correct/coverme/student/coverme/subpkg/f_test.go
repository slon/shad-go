package subpkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddOne(t *testing.T) {
	require.Equal(t, 1, AddOne(0))
}
