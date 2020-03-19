package subpkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddOne_2(t *testing.T) {
	require.Equal(t, 2, AddOne(1))
	require.Equal(t, 30, AddOne(29))
}

func TestAddTwo(t *testing.T) {
	require.Equal(t, 4, AddTwo(2))
}
