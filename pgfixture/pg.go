package pgfixture

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func Start(t *testing.T) string {
	pgconn, ok := os.LookupEnv("PGCONN")
	if ok {
		t.Logf("using external database: PGCONN=%s", pgconn)
		return pgconn
	}

	_, err := exec.LookPath("initdb")
	if err != nil {
		t.Fatalf("initdb binary not found; is postgres installed?")
	}

	_, err = exec.LookPath("postgres")
	if err != nil {
		t.Fatalf("postgres binary not found; is postgres installed?")
	}

	pgdata := t.TempDir()

	cmd := exec.Command("initdb", "-N", "-D", pgdata)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		t.Fatalf("initdb failed: %v", err)
	}

	pgrun := t.TempDir()

	cmd = exec.Command("postgres", "-D", pgdata, "-k", pgrun, "-F")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Start(); err != nil {
		t.Fatalf("postgres failed: %v", err)
	}

	finished := make(chan error, 1)
	go func() {
		finished <- cmd.Wait()
	}()

	select {
	case <-finished:
		t.Fatalf("postgres server terminated: %v", err)

	case <-time.After(time.Second / 2):
	}

	t.Cleanup(func() {
		_ = cmd.Process.Kill()
	})

	return fmt.Sprintf("host=%s database=postgres", pgrun)
}
