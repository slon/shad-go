package fileleak_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/fileleak"
)

type fakeT struct {
	failed  bool
	buffer  bytes.Buffer
	cleanup []func()
}

func (f *fakeT) Errorf(msg string, args ...interface{}) {
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

		ff, err := ioutil.TempFile("", "")
		require.NoError(t, err)
		f = ff
	})
}
