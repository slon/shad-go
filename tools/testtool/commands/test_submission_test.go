package commands

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// List directories in given directory.
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
	require.Nil(t, err)

	for _, dir := range testDirs {
		dir, err := filepath.Abs(dir)
		require.Nil(t, err)
		problem := path.Base(dir)
		t.Run(problem, func(t *testing.T) {
			studentRepo := path.Join(dir, "student")
			privateRepo := path.Join(dir, "private")
			testSubmission(studentRepo, privateRepo, problem)
		})
	}
}

func Test_testSubmission_incorrect(t *testing.T) {
	testDirs, err := listDirs("../testdata/submissions/incorrect")
	require.Nil(t, err)

	for _, dir := range testDirs {
		dir, err := filepath.Abs(dir)
		require.Nil(t, err)

		problem := path.Base(dir)
		t.Run(problem, func(t *testing.T) {
			if os.Getenv("BE_CRASHER") == "1" {
				studentRepo := path.Join(dir, "student")
				privateRepo := path.Join(dir, "private")
				testSubmission(studentRepo, privateRepo, problem)
				return
			}

			cmd := exec.Command(os.Args[0], "-v=0", "-test.run=Test_testSubmission_incorrect/"+problem)
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}

			t.Fatalf("process ran with err %v, want exit status != 0", err)
		})
	}
}
