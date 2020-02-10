package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func absPaths(files []string) []string {
	var abs []string
	for _, f := range files {
		absPath, _ := filepath.Abs("../testdata/list/" + f)
		abs = append(abs, absPath)
	}
	return abs
}

func TestListTestFiles(t *testing.T) {
	require.Equal(t,
		absPaths([]string{"sum/private_test.go", "sum/public_test.go"}),
		listTestFiles("../testdata/list"))
}

func TestProtectedFiles(t *testing.T) {
	require.Equal(t,
		absPaths([]string{"sum/dontchange.go"}),
		listProtectedFiles("../testdata/list"))
}

func TestPrivateFiles(t *testing.T) {
	require.Equal(t,
		absPaths([]string{"sum/private_test.go", "sum/solution.go"}),
		listPrivateFiles("../testdata/list"))
}

func TestListPackages(t *testing.T) {
	binaries, tests := listTestsAndBinaries("../testdata/pkgfind/task", []string{"-tags", "private"})

	assert.Equal(t, binaries, map[string]struct{}{
		"gitlab.com/slon/shad-go/task/cmd/tool":           {},
		"gitlab.com/slon/shad-go/task/cmd/tool_with_test": {},
	})

	assert.Equal(t, tests, map[string]struct{}{
		"gitlab.com/slon/shad-go/task/cmd/tool_with_test": {},
		"gitlab.com/slon/shad-go/task/pkg/a":              {},
		"gitlab.com/slon/shad-go/task/pkg/c":              {},
	})
}
