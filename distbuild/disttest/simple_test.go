package disttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var echoGraph = build.Graph{
	Jobs: []build.Job{
		{
			ID:   build.ID{'a'},
			Name: "echo",
			Cmds: []build.Cmd{
				{Exec: []string{"echo", "-n", "OK"}},
			},
		},
	},
}

func TestSingleCommand(t *testing.T) {
	env, cancel := newEnv(t)
	defer cancel()

	var recorder Recorder
	require.NoError(t, env.Client.Build(env.Ctx, echoGraph, &recorder))

	assert.Len(t, len(recorder.Jobs), 1)
	assert.Equal(t, &JobResult{Stdout: "OK", Code: new(int)}, recorder.Jobs[build.ID{'a'}])
}
