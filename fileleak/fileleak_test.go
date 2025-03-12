package fileleak_test

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/fileleak"
	"gitlab.com/slon/shad-go/tools/testtool"
)

type fakeT struct {
	failed  bool
	buffer  bytes.Buffer
	cleanup []func()
}

func (f *fakeT) Errorf(msg string, args ...any) {
	f.failed = true
	fmt.Fprintf(&f.buffer, msg, args...)
}

func (f *fakeT) Cleanup(cb func()) {
	f.cleanup = append([]func(){cb}, f.cleanup...)
}

func checkLeak(t *testing.T, shouldDetect bool, leaker func()) {
	var fake fakeT
	fileleak.VerifyNone(&fake)
	leaker()
	for _, cb := range fake.cleanup {
		cb()
	}

	switch {
	case shouldDetect && !fake.failed:
		t.Errorf("fileleak didn't detect a leak")
	case !shouldDetect && fake.failed:
		t.Errorf("fileleak detected a leak when there is none:\n%s", fake.buffer.String())
	}
}

func TestFileLeak_OpenFile(t *testing.T) {
	var f *os.File
	checkLeak(t, true, func() {
		var err error
		f, err = os.Open("/proc/self/exe")
		require.NoError(t, err)
	})

	if f != nil {
		_ = f.Close()
	}
}

func TestFileLeak_AlwaysOpenFile(t *testing.T) {
	f, err := os.Open("/proc/self/exe")
	require.NoError(t, err)
	defer f.Close()

	checkLeak(t, false, func() {})
}

func TestFileLeak_ReopenFile(t *testing.T) {
	f, err := os.Open("/proc/self/exe")
	require.NoError(t, err)
	defer f.Close()

	checkLeak(t, true, func() {
		_ = f.Close()

		ff, err := os.CreateTemp("", "")
		require.NoError(t, err)
		f = ff
	})
}

func TestFileLeak_PipeLeak(t *testing.T) {
	checkLeak(t, true, func() {
		f, _, err := os.Pipe()
		require.NoError(t, err)
		_ = f.Close()
	})
}

func TestFileLeak_PipeNoLeak(t *testing.T) {
	checkLeak(t, false, func() {
		f, w, err := os.Pipe()
		require.NoError(t, err)
		_, _ = f.Close(), w.Close()
	})
}

func TestFileLeak_SocketLeak(t *testing.T) {
	var conn net.Listener
	defer func() { conn.Close() }()

	checkLeak(t, true, func() {
		addr, err := testtool.GetFreePort()
		require.NoError(t, err)

		conn, err = net.Listen("tcp", fmt.Sprintf(":%s", addr))
		require.NoError(t, err)
	})
}

func TestFileLeak_SocketNoLeak(t *testing.T) {
	checkLeak(t, false, func() {
		addr, err := testtool.GetFreePort()
		require.NoError(t, err)

		conn, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
		require.NoError(t, err)
		_ = conn.Close()
	})
}

func TestFileLeak_DupLeak(t *testing.T) {
	var fd int
	defer syscall.Close(fd)

	checkLeak(t, true, func() {
		var err error
		fd, err = syscall.Dup(1)
		require.NoError(t, err)

	})
}

func TestFileLeak_DupNoLeak(t *testing.T) {
	checkLeak(t, false, func() {
		fd, err := syscall.Dup(1)
		require.NoError(t, err)
		_ = syscall.Close(fd)
	})
}
