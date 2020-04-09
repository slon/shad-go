package build

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTopSort(t *testing.T) {
	jobs := []Job{
		{
			ID:   ID{'a'},
			Deps: []ID{{'b'}},
		},
		{
			ID:   ID{'b'},
			Deps: []ID{{'c'}},
		},
		{
			ID: ID{'c'},
		},
	}

	sorted := TopSort(jobs)
	require.Equal(t, 3, len(sorted))
	require.Equal(t, ID{'c'}, sorted[0].ID)
	require.Equal(t, ID{'b'}, sorted[1].ID)
	require.Equal(t, ID{'a'}, sorted[2].ID)
}
