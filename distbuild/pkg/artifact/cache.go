package artifact

import (
	"errors"
	"net/http"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var (
	ErrNotFound    = errors.New("file not found")
	ErrWriteLocked = errors.New("file is locked for write")
	ErrReadLocked  = errors.New("file is locked for read")
)

type Cache struct{}

func NewCache(root string) (*Cache, error) {
	panic("implement me")
}

func (c *Cache) Range(artifactFn func(file build.ID) error) error {
	panic("implement me")
}

func (c *Cache) Remove(artifact build.ID) error {
	panic("implement me")
}

func (c *Cache) Create(artifact build.ID) (path string, abort, commit func(), err error) {
	panic("implement me")
}

func (c *Cache) Get(file build.ID) (path string, unlock func(), err error) {
	panic("implement me")
}

func NewHandler(c *Cache) http.Handler {
	panic("implement me")
}
