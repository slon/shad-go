package externalsort

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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
1`, `0
2`},
			out: `0
0
1
1
1
2
`,
		},
		{
			// Merge believes lines are read in sorted order.
			name: "single-unsorted-file",
			in: []string{`1
0`},
			out: `1
0
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			w := bufio.NewWriter(out)
			lw := NewWriter(w)

			var readers []LineReader
			for _, s := range tc.in {
				readers = append(readers, newStringReader(s))
			}

			err := Merge(lw, readers...)
			require.NoError(t, err)

			require.NoError(t, w.Flush())
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
			tmpDir, err := ioutil.TempDir("", fmt.Sprintf("sort%s-", d))
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			in, out := readTestCase(testCaseDir)
			in = copyFiles(t, in, tmpDir)

			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			require.NoError(t, Sort(w, in...))

			expected, err := ioutil.ReadFile(out)
			require.NoError(t, err)

			require.NoError(t, w.Flush())
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

func copyFiles(t *testing.T, in []string, dir string) []string {
	t.Helper()

	var ret []string
	for _, f := range in {
		ret = append(ret, copyFile(t, f, dir))
	}

	return ret
}

func copyFile(t *testing.T, f, dir string) string {
	t.Helper()

	data, err := ioutil.ReadFile(f)
	require.NoError(t, err)

	dst := path.Join(dir, path.Base(f))
	err = ioutil.WriteFile(dst, data, 0644)
	require.NoError(t, err)

	return dst
}
