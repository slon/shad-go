package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	files, err := listChangedFiles(".")
	require.NoError(t, err)
	require.NotEmpty(t, files)
}
