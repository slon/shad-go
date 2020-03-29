package filecache

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

func TestFileCache(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "filecache")
	require.NoError(t, err)

	cache, err := New(tmpDir)
	require.NoError(t, err)

	_, abort, err := cache.Write(build.ID{01})
	require.NoError(t, err)
	require.NoError(t, abort())

	_, _, err = cache.Get(build.ID{01})
	require.Truef(t, errors.Is(err, ErrNotFound), "real error: %v", err)

	f, _, err := cache.Write(build.ID{02})
	require.NoError(t, err)

	_, err = f.Write([]byte("foo bar"))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	path, unlock, err := cache.Get(build.ID{02})
	require.NoError(t, err)
	defer unlock()

	content, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, []byte("foo bar"), content)
}
