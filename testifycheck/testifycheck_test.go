package testifycheck

import (
	"os"
	"os/exec"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	debugOut, err := exec.Command("go", "env", "GOROOT").CombinedOutput()
	if err != nil {
		_, _ = os.Stderr.Write(debugOut)
	}

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "tests/...")
}
