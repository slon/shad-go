package filecache_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
)

type env struct {
	cache  *testCache
	server *httptest.Server
	client *filecache.Client
}

func newEnv(t *testing.T) *env {
	l := zaptest.NewLogger(t)
	mux := http.NewServeMux()

	cache := newCache(t)
	defer func() {
		if cache != nil {
			cache.cleanup()
		}
	}()

	handler := filecache.NewHandler(l, cache.Cache)
	handler.Register(mux)

	server := httptest.NewServer(mux)

	client := filecache.NewClient(l, server.URL)

	env := &env{
		cache:  cache,
		server: server,
		client: client,
	}

	cache = nil
	return env
}

func (e *env) stop() {
	e.server.Close()
	e.cache.cleanup()
}

func TestFileUpload(t *testing.T) {
	env := newEnv(t)
	defer env.stop()

	content := bytes.Repeat([]byte("foobar"), 1024*1024)

	tmpFilePath := filepath.Join(env.cache.tmpDir, "foo.txt")
	require.NoError(t, ioutil.WriteFile(tmpFilePath, content, 0666))

	ctx := context.Background()

	t.Run("UploadSingleFile", func(t *testing.T) {
		id := build.ID{0x01}

		require.NoError(t, env.client.Upload(ctx, id, tmpFilePath))

		path, unlock, err := env.cache.Get(id)
		require.NoError(t, err)
		defer unlock()

		actualContent, err := ioutil.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, content, actualContent)
	})

	t.Run("RepeatedUpload", func(t *testing.T) {
		id := build.ID{0x02}

		require.NoError(t, env.client.Upload(ctx, id, tmpFilePath))
		require.NoError(t, env.client.Upload(ctx, id, tmpFilePath))
	})

	t.Run("ConcurrentUpload", func(t *testing.T) {
		const (
			N = 10
			G = 10
		)

		for i := 0; i < N; i++ {
			var wg sync.WaitGroup
			wg.Add(G)

			id := build.ID{0x03, byte(i)}
			for j := 0; j < G; j++ {
				go func() {
					defer wg.Done()

					assert.NoError(t, env.client.Upload(ctx, id, tmpFilePath))
				}()
			}

			wg.Wait()
		}
	})
}

func TestFileDownload(t *testing.T) {
	env := newEnv(t)
	defer env.stop()

	localCache := newCache(t)
	defer localCache.cleanup()

	id := build.ID{0x01}

	w, abort, err := env.cache.Write(id)
	require.NoError(t, err)
	defer func() { _ = abort() }()

	_, err = w.Write([]byte("foobar"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	ctx := context.Background()
	require.NoError(t, env.client.Download(ctx, localCache.Cache, id))

	path, unlock, err := localCache.Get(id)
	require.NoError(t, err)
	defer unlock()

	content, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, []byte("foobar"), content)
}
