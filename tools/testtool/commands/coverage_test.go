package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getCoverageRequirements(t *testing.T) {
	r := getCoverageRequirements("../testdata/coverage/sum")
	require.True(t, r.Enabled)
	require.Equal(t, 90.0, r.Percent)
	require.Equal(t, []string{"."}, r.Packages)
}
