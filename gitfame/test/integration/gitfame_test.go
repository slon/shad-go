package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	testsDir := path.Join("./testdata", "good")
	files, err := ioutil.ReadDir(testsDir)
	require.NoError(t, err)

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		tc := ReadTestCase(t, filepath.Join(testsDir, f.Name()))

		t.Run(f.Name()+"/"+tc.Name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "gitfame-")
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(dir) }()

			args := []string{"--repository", dir}
			args = append(args, tc.Args...)

			Unbundle(t, filepath.Join(bundlesDir, tc.Bundle), dir)

			cmd := exec.Command(binary, args...)
			cmd.Stderr = ioutil.Discard

			output, err := cmd.Output()
			require.NoError(t, err)

			require.Equal(t, string(tc.Expected), string(output))
		})
	}
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
