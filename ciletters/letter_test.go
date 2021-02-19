package ciletters

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/tools/testtool"
)

type testCase struct {
	name         string
	notification Notification
	expected     string
}

func TestMakeLetter(t *testing.T) {
	const (
		testUser      = "gopher"
		gitlabGroupID = "go-spring-2021"
	)

	randomGitlabGroup := testtool.RandomName()
	randomGitlabProject := testtool.RandomName()
	randomBranch := testtool.RandomName()
	randomTestUser := testtool.RandomName()
	randomTriggerer := testtool.RandomName()
	randomPipelineID := randomInt64(t)
	randomJobIDS := []int64{randomInt64(t), randomInt64(t)}
	randomJobNames := []string{testtool.RandomName(), testtool.RandomName()}
	randomStages := []string{testtool.RandomName(), testtool.RandomName()}
	randomHash := testtool.RandomName()[:8]

	for _, tc := range []testCase{
		{
			name: "success",
			notification: Notification{
				Project: GitlabProject{
					GroupID: gitlabGroupID,
					ID:      testUser,
				},
				Branch: "master",
				Commit: Commit{
					Hash:    "2ff019bcb8f68d13d640e13351dad98edf7f1405",
					Message: "Solve sum.",
					Author:  testUser,
				},
				Pipeline: Pipeline{
					Status:      PipelineStatusOK,
					ID:          194555,
					TriggeredBy: testUser,
				},
			},
			expected: `Your pipeline #194555 passed!
    Project:      go-spring-2021/gopher
    Branch:       🌿 master
    Commit:       2ff019bc Solve sum.
    CommitAuthor: gopher`,
		},
		{
			name: "job-failed",
			notification: Notification{
				Project: GitlabProject{
					GroupID: gitlabGroupID,
					ID:      testUser,
				},
				Branch: "master",
				Commit: Commit{
					Hash:    "8967153e8aa7b270af6447dae594eb87bdae8791",
					Message: "Solve urlfetch.",
					Author:  testUser,
				},
				Pipeline: Pipeline{
					Status:      PipelineStatusFailed,
					ID:          194613,
					TriggeredBy: testUser,
					FailedJobs: []Job{
						{
							ID:    202538,
							Name:  "grade",
							Stage: "test",
							RunnerLog: `$ testtool grade
testtool: detected change in tasks [sum]
testtool: skipping task sum: not released yet
testtool: testing task sum
testtool: testing submission in /tmp/sum-281145206
testtool: copying student repo
testtool: copying tests
testtool: copying !change files
testtool: copying testdata directory
testtool: copying go.mod, go.sum and .golangci.yml
testtool: running tests
testtool: > go test -mod readonly -tags private -c -o /tmp/bincache730817117/5d83984f885e61c1 gitlab.com/slon/shad-go/sum
--- FAIL: TestSum (0.00s)
    sum_test.go:19: 2 + 2 == 0 != 4
    sum_test.go:19: 9223372036854775807 + 1 == 0 != -9223372036854775808
FAIL
testtool: task sum failed: test failed: exit status 1
some tasks failed
ERROR: Job failed: exit code 1`,
						},
					},
				},
			},
			expected: `Your pipeline #194613 has failed!
    Project:      go-spring-2021/gopher
    Branch:       🌿 master
    Commit:       8967153e Solve urlfetch.
    CommitAuthor: gopher
        Stage: test, Job grade
            testtool: copying go.mod, go.sum and .golangci.yml
            testtool: running tests
            testtool: > go test -mod readonly -tags private -c -o /tmp/bincache730817117/5d83984f885e61c1 gitlab.com/slon/shad-go/sum
            --- FAIL: TestSum (0.00s)
                sum_test.go:19: 2 + 2 == 0 != 4
                sum_test.go:19: 9223372036854775807 + 1 == 0 != -9223372036854775808
            FAIL
            testtool: task sum failed: test failed: exit status 1
            some tasks failed
            ERROR: Job failed: exit code 1
`,
		},
		{
			name: "multiple-jobs-failed",
			notification: Notification{
				Project: GitlabProject{
					GroupID: randomGitlabGroup,
					ID:      randomGitlabProject,
				},
				Branch: randomBranch,
				Commit: Commit{
					Hash:    randomHash,
					Message: "Solve digitalclock.",
					Author:  randomTestUser,
				},
				Pipeline: Pipeline{
					Status:      PipelineStatusFailed,
					ID:          randomPipelineID,
					TriggeredBy: randomTriggerer,
					FailedJobs: []Job{
						{
							ID:    randomJobIDS[0],
							Name:  randomJobNames[0],
							Stage: randomStages[0],
							RunnerLog: `$ testtool grade
testtool: detected change in tasks [sum]
testtool: skipping task sum: not released yet
testtool: testing task sum
testtool: testing submission in /tmp/sum-281145206
testtool: copying student repo
testtool: copying tests
testtool: copying !change files
testtool: copying testdata directory
testtool: copying go.mod, go.sum and .golangci.yml
testtool: running tests
testtool: > go test -mod readonly -tags private -c -o /tmp/bincache730817117/5d83984f885e61c1 gitlab.com/slon/shad-go/sum
--- FAIL: TestSum (0.00s)
    sum_test.go:19: 2 + 2 == 0 != 4
    sum_test.go:19: 9223372036854775807 + 1 == 0 != -9223372036854775808
FAIL
testtool: task sum failed: test failed: exit status 1
some tasks failed
ERROR: Job failed: exit code 1`,
						},
						{
							ID:    randomJobIDS[1],
							Name:  randomJobNames[1],
							Stage: randomStages[1],
							RunnerLog: `--- FAIL: TestSum (0.00s)
    sum_test.go:19: 2 + 2 == 0 != 4
    sum_test.go:19: 9223372036854775807 + 1 == 0 != -9223372036854775808
FAIL
testtool: task sum failed: test failed: exit status 1
some tasks failed
ERROR: Job failed: exit code 1`,
						},
					},
				},
			},
			expected: fmt.Sprintf(`Your pipeline #%d has failed!
    Project:      %v/%v
    Branch:       🌿 %v
    Commit:       %v Solve digitalclock.
    CommitAuthor: %v
        Stage: %v, Job %v
            testtool: copying go.mod, go.sum and .golangci.yml
            testtool: running tests
            testtool: > go test -mod readonly -tags private -c -o /tmp/bincache730817117/5d83984f885e61c1 gitlab.com/slon/shad-go/sum
            --- FAIL: TestSum (0.00s)
                sum_test.go:19: 2 + 2 == 0 != 4
                sum_test.go:19: 9223372036854775807 + 1 == 0 != -9223372036854775808
            FAIL
            testtool: task sum failed: test failed: exit status 1
            some tasks failed
            ERROR: Job failed: exit code 1

        Stage: %v, Job %v
            --- FAIL: TestSum (0.00s)
                sum_test.go:19: 2 + 2 == 0 != 4
                sum_test.go:19: 9223372036854775807 + 1 == 0 != -9223372036854775808
            FAIL
            testtool: task sum failed: test failed: exit status 1
            some tasks failed
            ERROR: Job failed: exit code 1
`, randomPipelineID, randomGitlabGroup, randomGitlabProject, randomBranch, randomHash, randomTestUser,
				randomStages[0], randomJobNames[0], randomStages[1], randomJobNames[1]),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			letter, err := MakeLetter(&tc.notification)
			require.NoError(t, err)
			require.Equal(t, tc.expected, letter)
		})
	}
}

func randomInt64(t *testing.T) int64 {
	t.Helper()

	nBig, err := rand.Int(rand.Reader, big.NewInt(1e6))
	require.NoError(t, err)

	return nBig.Int64()
}
