package shopfront_test

import (
	"os"
	"os/exec"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/tools/testtool"
)

type testingTB interface {
	Logf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Cleanup(func())
}

func StartRedis(tb testingTB) string {
	if redis, ok := os.LookupEnv("REDIS"); ok {
		tb.Logf("using external redis server; REDIS=%s", redis)
		return redis
	}

	port, err := testtool.GetFreePort()
	require.NoError(tb, err)

	_, err = exec.LookPath("redis-server")
	if err != nil {
		tb.Fatalf("redis-server binary is not found; is redis installed?")
	}

	cmd := exec.Command("redis-server", "--port", port, "--save", "", "--appendonly", "no")
	cmd.Stderr = os.Stderr

	require.NoError(tb, cmd.Start())

	finished := make(chan error, 1)
	go func() {
		finished <- cmd.Wait()
	}()

	select {
	case err := <-finished:
		tb.Fatalf("redis server terminated: %v", err)

	case <-time.After(time.Second / 2):
	}

	tb.Cleanup(func() {
		_ = cmd.Process.Kill()
	})

	return "localhost:" + port
}
