package disttest

import (
	"io/ioutil"
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
				{Exec: []string{"echo", "OK"}},
			},
		},
	},
}

func TestSingleCommand(t *testing.T) {
	env, cancel := newEnv(t)
	defer cancel()

	recorder := NewRecorder()
	require.NoError(t, env.Client.Build(env.Ctx, echoGraph, recorder))

	assert.Len(t, recorder.Jobs, 1)
	assert.Equal(t, &JobResult{Stdout: "OK\n", Code: new(int)}, recorder.Jobs[build.ID{'a'}])
}

func TestJobCaching(t *testing.T) {
	env, cancel := newEnv(t)
	defer cancel()

	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)

	graph := build.Graph{
		Jobs: []build.Job{
			{
				ID:   build.ID{'a'},
				Name: "echo",
				Cmds: []build.Cmd{
					{CatTemplate: "OK\n", CatOutput: tmpFile.Name()}, // No-hermetic, for testing purposes.
					{Exec: []string{"echo", "OK"}},
				},
			},
		},
	}

	recorder := NewRecorder()
	require.NoError(t, env.Client.Build(env.Ctx, graph, recorder))

	assert.Len(t, recorder.Jobs, 1)
	assert.Equal(t, &JobResult{Stdout: "OK\n", Code: new(int)}, recorder.Jobs[build.ID{'a'}])

	// Second build must get results from cache.
	require.NoError(t, env.Client.Build(env.Ctx, graph, NewRecorder()))

	output, err := ioutil.ReadAll(tmpFile)
	require.NoError(t, err)
	require.Equal(t, []byte("OK\n"), output)
}

func TestSourceFiles(t *testing.T) {

}
