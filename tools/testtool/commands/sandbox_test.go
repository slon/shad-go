package commands

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSandbox(t *testing.T) {
	var cmd exec.Cmd

	require.NoError(t, sandbox(&cmd))
	require.True(t, cmd.SysProcAttr.Credential.Uid > 0)
	require.True(t, cmd.SysProcAttr.Credential.Gid > 0)
}
