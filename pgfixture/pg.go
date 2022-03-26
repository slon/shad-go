package pgfixture

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func lookPath(t *testing.T, name string) string {
	t.Helper()

	path, err := exec.LookPath(name)
	if err == nil {
		return path
	}

	const ubuntuPostgres = "/usr/lib/postgresql"

	if dirs, err := ioutil.ReadDir(ubuntuPostgres); err == nil {
		for _, d := range dirs {
			path := filepath.Join(ubuntuPostgres, d.Name(), "bin", name)

			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	t.Fatalf("%s binary not found; is postgres installed?", name)
	return ""
}

func Start(t *testing.T) string {
	pgconn, ok := os.LookupEnv("PGCONN")
	if ok {
		t.Logf("using external database: PGCONN=%s", pgconn)
		return pgconn
	}

	pgdata := t.TempDir()

	cmd := exec.Command(lookPath(t, "initdb"), "-N", "-D", pgdata)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		t.Fatalf("initdb failed: %v", err)
	}

	pgrun := t.TempDir()

	cmd = exec.Command(lookPath(t, "postgres"), "-D", pgdata, "-k", pgrun, "-F")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		t.Fatalf("postgres failed: %v", err)
	}

	finished := make(chan error, 1)
	go func() {
		finished <- cmd.Wait()
	}()

	select {
	case err := <-finished:
		t.Fatalf("postgres server terminated: %v", err)

	case <-time.After(time.Second / 2):
	}

	t.Cleanup(func() {
		_ = cmd.Process.Kill()
	})

	return fmt.Sprintf("host=%s database=postgres", pgrun)
}
