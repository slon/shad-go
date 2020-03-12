package externalsort

import (
	"bytes"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

func TestMerge(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "simple",
			in: []string{`0`, `1
1
1`},
			out: `0
1
1
1`,
		},
		{
			// Merge believes lines are read in sorted order.
			name: "single-unsorted-file",
			in: []string{`1
0`},
			out: `1
0`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			w := NewWriterFlusher(out)

			var readers []LineReader
			for _, s := range tc.in {
				readers = append(readers, newStringReader(s))
			}

			err := Merge(w, readers...)
			require.NoError(t, err)

			require.Equal(t, tc.out, out.String())
		})
	}
}

func TestSort_fileNotFound(t *testing.T) {
	var buf bytes.Buffer
	err := Sort(&buf, testtool.RandomName())
	require.Error(t, err)
}

func TestSort(t *testing.T) {
	testDir := path.Join("./testdata", "sort")

	readTestCase := func(dir string) (in []string, out string) {
		files, err := ioutil.ReadDir(dir)
		require.NoError(t, err)

		for _, f := range files {
			if strings.HasPrefix(f.Name(), "in") {
				in = append(in, path.Join(dir, f.Name()))
			}
			if f.Name() == "out.txt" {
				out = path.Join(dir, f.Name())
			}
		}

		return
	}

	for _, d := range listDirs(t, testDir) {
		testCaseDir := path.Join(testDir, d)

		t.Run(d, func(t *testing.T) {
			in, out := readTestCase(testCaseDir)

			var buf bytes.Buffer
			err := Sort(&buf, in...)
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(out)
			require.NoError(t, err)

			require.Equal(t, string(expected), buf.String())
		})
	}
}

func listDirs(t *testing.T, dir string) []string {
	t.Helper()

	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err)

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	return dirs
}
