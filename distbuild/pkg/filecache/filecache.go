package filecache

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var (
	ErrNotFound    = errors.New("file not found")
	ErrWriteLocked = errors.New("file is locked for write")
	ErrReadLocked  = errors.New("file is locked for read")
)

type Cache struct {
	rootDir string
}

func New(rootDir string) (*Cache, error) {
	if err := os.MkdirAll(rootDir, 0777); err != nil {
		return nil, fmt.Errorf("error creating filecache: %w", err)
	}

	c := &Cache{rootDir: rootDir}
	return c, nil
}

func (c *Cache) Range(fileFn func(file build.ID) error) error {
	panic("implement me")
}

func (c *Cache) Remove(file build.ID) error {
	panic("implement me")
}

func (c *Cache) Write(file build.ID) (w io.WriteCloser, abort func(), err error) {
	panic("implement me")
}

func (c *Cache) Get(file build.ID) (path string, unlock func(), err error) {
	panic("implement me")
}

func NewHandler(c *Cache) http.Handler {
	panic("implement me")
}
