package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeadlines(t *testing.T) {
	d, err := loadDeadlines("../../../.deadlines.yml")
	require.NoError(t, err)
	require.NotEmpty(t, d)

	_, sum := d.FindTask("sum")
	require.NotNil(t, sum)
	require.Equal(t, "sum", sum.Name)
}

func TestDetectChange(t *testing.T) {
	for _, tc := range []struct {
		name         string
		deadlines    string
		changedFiles []string
		changedTasks []string
	}{
		{
			name:         "sum", // Original deadlines file with sum task.
			deadlines:    "../../../.deadlines.yml",
			changedFiles: []string{"sum/sum.go", "testtool/foo.go", "README.md"},
			changedTasks: []string{"sum"},
		},
		{
			name:      "tarstreamtest", // Deadlines file with tarstreamtest task.
			deadlines: "../testdata/deadlines/.deadlines.yml",
			changedFiles: []string{
				"sum/sum.go",
				"testtool/foo.go",
				"distbuild/pkg/tarstream/stream.go",
				"README.md",
			},
			changedTasks: []string{"sum", "tarstreamtest"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := loadDeadlines(tc.deadlines)
			require.NoError(t, err)

			changed := findChangedTasks(d, tc.changedFiles)
			require.Equal(t, tc.changedTasks, changed)
		})
	}
}
