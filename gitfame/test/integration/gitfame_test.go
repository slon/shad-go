package integration

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"gitlab.com/slon/shad-go/tools/testtool"
)

const importPath = "gitlab.com/slon/shad-go/gitfame/cmd/gitfame"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

func TestGitFame(t *testing.T) {
	binary, err := binCache.GetBinary(importPath)
	require.NoError(t, err)

	bundlesDir := path.Join("./testdata", "bundles")
	testsDir := path.Join("./testdata", "tests")
	testDirs := ListTestDirs(t, testsDir)

	for _, dir := range testDirs {
		tc := ReadTestCase(t, filepath.Join(testsDir, dir))

		t.Run(dir+"/"+tc.Name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "gitfame-")
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(dir) }()

			args := []string{"--repository", dir}
			args = append(args, tc.Args...)

			Unbundle(t, filepath.Join(bundlesDir, tc.Bundle), dir)
			headRef := GetHEADRef(t, dir)

			cmd := exec.Command(binary, args...)
			cmd.Stderr = os.Stderr

			output, err := cmd.Output()
			if !tc.Error {
				require.NoError(t, err)
				CompareResults(t, tc.Expected, output, tc.Format)
			} else {
				require.Error(t, err)
				_, ok := err.(*exec.ExitError)
				require.True(t, ok)
			}

			newHEADRef := GetHEADRef(t, dir)
			require.Equal(t, headRef, newHEADRef)
		})
	}
}

func ListTestDirs(t *testing.T, path string) []string {
	t.Helper()

	files, err := ioutil.ReadDir(path)
	require.NoError(t, err)

	var names []string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		names = append(names, f.Name())
	}

	toInt := func(name string) int {
		i, err := strconv.Atoi(name)
		require.NoError(t, err)
		return i
	}

	sort.Slice(names, func(i, j int) bool {
		return toInt(names[i]) < toInt(names[j])
	})

	return names
}

type TestCase struct {
	*TestDescription
	Expected []byte
}

func ReadTestCase(t *testing.T, path string) *TestCase {
	t.Helper()

	desc := ReadTestDescription(t, path)

	expected, err := ioutil.ReadFile(filepath.Join(path, "expected.out"))
	require.NoError(t, err)

	return &TestCase{TestDescription: desc, Expected: expected}
}

type TestDescription struct {
	Name   string   `yaml:"name"`
	Args   []string `yaml:"args"`
	Bundle string   `yaml:"bundle"`
	Error  bool     `yaml:"error"`
	Format string   `yaml:"format,omitempty"`
}

func ReadTestDescription(t *testing.T, path string) *TestDescription {
	t.Helper()

	data, err := ioutil.ReadFile(filepath.Join(path, "description.yaml"))
	require.NoError(t, err)

	var desc TestDescription
	require.NoError(t, yaml.Unmarshal(data, &desc))

	return &desc
}

func Unbundle(t *testing.T, src, dst string) {
	t.Helper()

	cmd := exec.Command("git", "clone", src, dst)
	require.NoError(t, cmd.Run())
}

func CompareResults(t *testing.T, expected, actual []byte, format string) {
	t.Helper()

	switch format {
	case "json":
		if len(expected) == 0 {
			require.Empty(t, string(actual))
		} else {
			require.JSONEq(t, string(expected), string(actual))
		}
	case "json-lines":
		if len(expected) == 0 {
			require.Empty(t, string(actual))
		} else {
			CompareJSONLines(t, expected, actual)
		}
	default:
		require.Equal(t, string(expected), string(actual))
	}
}

func CompareJSONLines(t *testing.T, expected, actual []byte) {
	t.Helper()

	expectedLines := ParseJSONLines(expected)
	actualLines := ParseJSONLines(actual)
	require.Equal(t, len(expectedLines), len(actualLines))

	for i, l := range expectedLines {
		require.JSONEq(t, string(l), string(actualLines[i]))
	}
}

func ParseJSONLines(data []byte) [][]byte {
	return bytes.Split(bytes.TrimSpace(data), []byte("\n"))
}

func GetHEADRef(t *testing.T, path string) string {
	t.Helper()

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path

	out, err := cmd.Output()
	require.NoError(t, err)

	return string(out)
}
