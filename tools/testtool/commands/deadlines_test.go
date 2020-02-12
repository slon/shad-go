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
	d, err := loadDeadlines("../../../.deadlines.yml")
	require.NoError(t, err)

	changed := findChangedTasks(d, []string{"sum/sum.go", "testtool/foo.go", "README.md"})
	require.Equal(t, []string{"sum"}, changed)
}
