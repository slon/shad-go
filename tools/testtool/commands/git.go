package commands

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func listChangedFiles(gitPath string) ([]string, error) {
	var gitOutput bytes.Buffer

	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD")
	cmd.Dir = gitPath
	cmd.Stdout = &gitOutput
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return strings.Split(gitOutput.String(), "\n"), nil
}
