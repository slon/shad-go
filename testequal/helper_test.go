package testequal

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelper(t *testing.T) {
	if os.Getenv("FAIL_ASSERTIONS") == "1" {
		AssertEqual(t, 1, 2, "%d must be equal to %d", 1, 2)
		AssertNotEqual(t, 1, 1, "1 != 1")
		RequireEqual(t, 1, 2)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.v", "-test.run=TestHelper")
	cmd.Env = append(os.Environ(), "FAIL_ASSERTIONS=1")
	var buf bytes.Buffer
	cmd.Stdout = &buf

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		require.Contains(t, buf.String(), "helper_test.go:14")
		require.Contains(t, buf.String(), "helper_test.go:15")
		require.Contains(t, buf.String(), "helper_test.go:16")
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
