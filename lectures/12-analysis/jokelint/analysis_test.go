package jokelint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalysis(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "tests/...")
}
