// +build !change

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

const wordcountImportPath = "gitlab.com/slon/shad-go/wordcount"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func TestWordCount(t *testing.T) {
	binary, err := binCache.GetBinary(wordcountImportPath)
	require.NoError(t, err)

	type counts map[string]int64
	type files []string

	for _, tc := range []struct {
		name     string
		files    files
		expected map[string]int64
	}{
		{
			name:     "empty",
			files:    files{``},
			expected: make(counts),
		},
		{
			name: "simple",
			files: files{
				`a
a
b
с


a
b`,
			},
			expected: counts{"a": 3, "b": 2, "": 2},
		},
		{
			name: "multiple-files",
			files: files{
				`a
a`,
				`a
b`,
				`b`,
			},
			expected: counts{"a": 3, "b": 2},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp directory.
			testDir, err := ioutil.TempDir("", "wordcount-testdata-")
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(testDir) }()

			// Create test files in temp directory.
			var files []string
			for _, f := range tc.files {
				file := path.Join(testDir, testtool.RandomName())
				err = ioutil.WriteFile(file, []byte(f), 0644)
				require.NoError(t, err)
				files = append(files, file)
			}

			// Run wordcount executable.
			cmd := exec.Command(binary, files...)
			cmd.Stderr = os.Stderr

			output, err := cmd.Output()
			require.NoError(t, err)

			// Parse output and compare with an expected one.
			counts, err := parseStdout(output)
			require.NoError(t, err)

			require.True(t, reflect.DeepEqual(tc.expected, counts),
				fmt.Sprintf("expected: %v; got: %v", tc.expected, counts))
		})
	}
}

// parseStdout parses wordcount's output of the ['<COUNT>\t<LINE>'] format.
func parseStdout(data []byte) (map[string]int64, error) {
	counts := make(map[string]int64)

	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected line format: %s", parts)
		}
		c, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse line count: %w", err)
		}
		counts[parts[1]] = c
	}

	return counts, nil
}
