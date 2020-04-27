package disttest

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var singleWorkerConfig = &Config{WorkerCount: 1}

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
	env, cancel := newEnv(t, singleWorkerConfig)
	defer cancel()

	recorder := NewRecorder()
	require.NoError(t, env.Client.Build(env.Ctx, echoGraph, recorder))

	assert.Len(t, recorder.Jobs, 1)
	assert.Equal(t, &JobResult{Stdout: "OK\n", Code: new(int)}, recorder.Jobs[build.ID{'a'}])
}

func TestJobCaching(t *testing.T) {
	env, cancel := newEnv(t, singleWorkerConfig)
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

	require.NoError(t, ioutil.WriteFile(tmpFile.Name(), []byte("NOTOK\n"), 0666))

	// Second build must get results from cache.
	require.NoError(t, env.Client.Build(env.Ctx, graph, NewRecorder()))

	output, err := ioutil.ReadAll(tmpFile)
	require.NoError(t, err)
	require.Equal(t, []byte("NOTOK\n"), output)
}

var sourceFilesGraph = build.Graph{
	SourceFiles: map[build.ID]string{
		{'a'}: "a.txt",
		{'c'}: "b/c.txt",
	},
	Jobs: []build.Job{
		{
			ID:   build.ID{'a'},
			Name: "echo",
			Cmds: []build.Cmd{
				{Exec: []string{"cat", "{{.SourceDir}}/a.txt"}},
				{Exec: []string{"bash", "-c", "cat {{.SourceDir}}/b/c.txt > /dev/stderr"}},
			},
			Inputs: []string{
				"a.txt",
				"b/c.txt",
			},
		},
	},
}

func TestSourceFiles(t *testing.T) {
	env, cancel := newEnv(t, singleWorkerConfig)
	defer cancel()

	recorder := NewRecorder()
	require.NoError(t, env.Client.Build(env.Ctx, sourceFilesGraph, recorder))

	assert.Len(t, recorder.Jobs, 1)
	assert.Equal(t, &JobResult{Stdout: "foo", Stderr: "bar", Code: new(int)}, recorder.Jobs[build.ID{'a'}])
}

var artifactTransferGraph = build.Graph{
	Jobs: []build.Job{
		{
			ID:   build.ID{'a'},
			Name: "write",
			Cmds: []build.Cmd{
				{CatTemplate: "OK", CatOutput: "{{.OutputDir}}/out.txt"},
			},
		},
		{
			ID:   build.ID{'b'},
			Name: "cat",
			Cmds: []build.Cmd{
				{Exec: []string{"cat", fmt.Sprintf("{{index .Deps %q}}/out.txt", build.ID{'a'})}},
			},
			Deps: []build.ID{{'a'}},
		},
	},
}

func TestArtifactTransferBetweenJobs(t *testing.T) {
	env, cancel := newEnv(t, singleWorkerConfig)
	defer cancel()

	recorder := NewRecorder()
	require.NoError(t, env.Client.Build(env.Ctx, artifactTransferGraph, recorder))

	assert.Len(t, recorder.Jobs, 2)
	assert.Equal(t, &JobResult{Stdout: "OK", Code: new(int)}, recorder.Jobs[build.ID{'b'}])
}
