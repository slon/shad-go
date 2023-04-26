package disttest

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var threeWorkerConfig = &Config{WorkerCount: 3}

func TestArtifactTransferBetweenWorkers(t *testing.T) {
	env, cancel := newEnv(t, threeWorkerConfig)
	defer cancel()

	baseJob := build.Job{
		ID:   build.ID{'a'},
		Name: "write",
		Cmds: []build.Cmd{
			{CatTemplate: "OK", CatOutput: "{{.OutputDir}}/out.txt"},
		},
	}

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		depJobID := build.ID{'b', byte(i)}
		depJob := build.Job{
			ID:   depJobID,
			Name: "cat",
			Cmds: []build.Cmd{
				{Exec: []string{"cat", fmt.Sprintf("{{index .Deps %q}}/out.txt", build.ID{'a'})}},
				{Exec: []string{"sleep", "1"}, Environ: os.Environ()}, // DepTimeout is 100ms.
			},
			Deps: []build.ID{{'a'}},
		}

		graph := build.Graph{Jobs: []build.Job{baseJob, depJob}}
		go func() {
			defer wg.Done()

			recorder := NewRecorder()
			if !assert.NoError(t, env.Client.Build(env.Ctx, graph, recorder)) {
				return
			}

			assert.Len(t, recorder.Jobs, 2)
			assert.Equal(t, &JobResult{Stdout: "OK", Code: new(int)}, recorder.Jobs[depJobID])
		}()
	}

	wg.Wait()

	for _, cache := range env.WorkerCache {
		_, unlock, err := cache.Get(baseJob.ID)
		require.NoError(t, err)
		defer unlock()
	}
}
