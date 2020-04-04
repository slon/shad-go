package artifact_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

func TestArtifactTransfer(t *testing.T) {
	remoteCache := newTestCache(t)
	defer remoteCache.cleanup()
	localCache := newTestCache(t)
	defer localCache.cleanup()

	id := build.ID{0x01}

	dir, commit, _, err := remoteCache.Create(id)
	require.NoError(t, err)
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "a.txt"), []byte("foobar"), 0777))
	require.NoError(t, commit())

	l := zaptest.NewLogger(t)

	h := artifact.NewHandler(l, remoteCache.Cache)
	mux := http.NewServeMux()
	h.Register(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	ctx := context.Background()
	require.NoError(t, artifact.Download(ctx, server.URL, localCache.Cache, id))

	dir, unlock, err := localCache.Get(id)
	require.NoError(t, err)
	defer unlock()

	content, err := ioutil.ReadFile(filepath.Join(dir, "a.txt"))
	require.NoError(t, err)
	require.Equal(t, []byte("foobar"), content)

	err = artifact.Download(ctx, server.URL, localCache.Cache, build.ID{0x02})
	require.Error(t, err)
}
