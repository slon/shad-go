package commands

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// listDirs lists directories in given directory.
func listDirs(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, path.Join(dir, f.Name()))
		}
	}

	return dirs, nil
}

func Test_testSubmission_correct(t *testing.T) {
	testDirs, err := listDirs("../testdata/submissions/correct")
	require.NoError(t, err)

	for _, dir := range testDirs {
		absDir, err := filepath.Abs(dir)
		require.NoError(t, err)
		problem := path.Base(absDir)
		t.Run(problem, func(t *testing.T) {
			studentRepo := path.Join(absDir, "student")
			privateRepo := path.Join(absDir, "private")

			require.NoError(t, testSubmission(studentRepo, privateRepo, problem))
		})
	}
}

func Test_testSubmission_incorrect(t *testing.T) {
	testDirs, err := listDirs("../testdata/submissions/incorrect")
	require.NoError(t, err)

	for _, dir := range testDirs {
		absDir, err := filepath.Abs(dir)
		require.NoError(t, err)

		problem := path.Base(absDir)
		t.Run(problem, func(t *testing.T) {
			studentRepo := path.Join(absDir, "student")
			privateRepo := path.Join(absDir, "private")

			err := testSubmission(studentRepo, privateRepo, problem)
			require.Error(t, err)

			if problem == "brokentest" {
				var testFailedErr *TestFailedError
				require.True(t, errors.As(err, &testFailedErr))
			}
		})
	}
}
