package artifact

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

func TestCache(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	c, err := NewCache(tmpDir)
	require.NoError(t, err)

	idA := build.ID{'a'}

	path, commit, _, err := c.Create(idA)
	require.NoError(t, err)

	_, _, _, err = c.Create(idA)
	require.Equal(t, ErrWriteLocked, err)

	_, err = os.Create(filepath.Join(path, "a.txt"))
	require.NoError(t, err)

	require.NoError(t, commit())

	path, unlock, err := c.Get(idA)
	require.NoError(t, err)
	defer unlock()

	_, err = os.Stat(filepath.Join(path, "a.txt"))
	require.NoError(t, err)

	require.Equal(t, ErrReadLocked, c.Remove(idA))

	idB := build.ID{'b'}
	_, _, err = c.Get(idB)
	require.Equal(t, ErrNotFound, err)

	require.NoError(t, c.Range(func(artifact build.ID) error {
		require.Equal(t, idA, artifact)
		return nil
	}))
}

func TestAbortWrite(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	c, err := NewCache(tmpDir)
	require.NoError(t, err)

	idA := build.ID{'a'}

	_, _, abort, err := c.Create(idA)
	require.NoError(t, err)
	require.NoError(t, abort())

	_, _, err = c.Get(idA)
	require.Equal(t, ErrNotFound, err)
}
