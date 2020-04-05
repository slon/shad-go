// +build !solution

package artifact

import (
	"errors"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var (
	ErrNotFound    = errors.New("artifact not found")
	ErrExists      = errors.New("artifact exists")
	ErrWriteLocked = errors.New("artifact is locked for write")
	ErrReadLocked  = errors.New("artifact is locked for read")
)

type Cache struct {
}

func NewCache(root string) (*Cache, error) {
	panic("implement me")
}

func (c *Cache) Range(artifactFn func(artifact build.ID) error) error {
	panic("implement me")
}

func (c *Cache) Remove(artifact build.ID) error {
	panic("implement me")
}

func (c *Cache) Create(artifact build.ID) (path string, commit, abort func() error, err error) {
	panic("implement me")
}

func (c *Cache) Get(artifact build.ID) (path string, unlock func(), err error) {
	panic("implement me")
}
