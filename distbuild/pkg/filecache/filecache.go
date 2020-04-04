package filecache

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var (
	ErrNotFound    = errors.New("file not found")
	ErrExists      = errors.New("file exists")
	ErrWriteLocked = errors.New("file is locked for write")
	ErrReadLocked  = errors.New("file is locked for read")
)

const fileName = "file"

func convertErr(err error) error {
	switch {
	case errors.Is(err, artifact.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, artifact.ErrExists):
		return ErrExists
	case errors.Is(err, artifact.ErrWriteLocked):
		return ErrWriteLocked
	case errors.Is(err, artifact.ErrReadLocked):
		return ErrReadLocked
	default:
		return err
	}
}

type Cache struct {
	cache *artifact.Cache
}

func New(rootDir string) (*Cache, error) {
	cache, err := artifact.NewCache(rootDir)
	if err != nil {
		return nil, err
	}

	c := &Cache{cache: cache}
	return c, nil
}

func (c *Cache) Range(fileFn func(file build.ID) error) error {
	return c.cache.Range(fileFn)
}

func (c *Cache) Remove(file build.ID) error {
	return convertErr(c.cache.Remove(file))
}

type fileWriter struct {
	f      *os.File
	commit func() error
}

func (f *fileWriter) Write(p []byte) (int, error) {
	return f.f.Write(p)
}

func (f *fileWriter) Close() error {
	closeErr := f.f.Close()
	commitErr := f.commit()

	if closeErr != nil {
		return closeErr
	}

	return commitErr
}

func (c *Cache) Write(file build.ID) (w io.WriteCloser, abort func() error, err error) {
	path, commit, abortDir, err := c.cache.Create(file)
	if err != nil {
		err = convertErr(err)
		return
	}

	f, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		_ = abort()
		return
	}

	w = &fileWriter{f: f, commit: commit}
	abort = func() error {
		closeErr := f.Close()
		abortErr := abortDir()

		if closeErr != nil {
			return closeErr
		}

		return abortErr
	}
	return
}

func (c *Cache) Get(file build.ID) (path string, unlock func(), err error) {
	root, unlock, err := c.cache.Get(file)
	path = filepath.Join(root, fileName)
	err = convertErr(err)
	return
}
