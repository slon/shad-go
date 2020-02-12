package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReport(t *testing.T) {
	if testingToken == "" {
		t.Skip("token is missing")
	}

	require.NoError(t, reportTestResults(testingToken, "sum", "1", false))
}
