package disttest

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	startTime := time.Now()

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

	testDuration := time.Since(startTime)
	assert.True(t, testDuration < time.Second*5/2, "test duration should be less than 2.5 seconds")
}
