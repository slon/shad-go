package artifact

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

var (
	ErrNotFound    = errors.New("file not found")
	ErrWriteLocked = errors.New("file is locked for write")
	ErrReadLocked  = errors.New("file is locked for read")
)

type Cache struct {
	tmpDir   string
	cacheDir string

	mu          sync.Mutex
	writeLocked map[build.ID]struct{}
	readLocked  map[build.ID]int
}

func NewCache(root string) (*Cache, error) {
	tmpDir := filepath.Join(root, "tmp")

	if err := os.RemoveAll(tmpDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(root, "c")
	if err := os.MkdirAll(cacheDir, 0777); err != nil {
		return nil, err
	}

	for i := 0; i < 256; i++ {
		d := hex.EncodeToString([]byte{uint8(i)})
		if err := os.MkdirAll(filepath.Join(cacheDir, d), 0777); err != nil {
			return nil, err
		}
	}

	return &Cache{
		tmpDir:      tmpDir,
		cacheDir:    cacheDir,
		writeLocked: make(map[build.ID]struct{}),
		readLocked:  make(map[build.ID]int),
	}, nil
}

func (c *Cache) readLock(id build.ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.writeLocked[id]; ok {
		return ErrWriteLocked
	}

	c.readLocked[id] += 1
	return nil
}

func (c *Cache) readUnlock(id build.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.readLocked[id] -= 1
	if c.readLocked[id] == 0 {
		delete(c.readLocked, id)
	}
}

func (c *Cache) writeLock(id build.ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.writeLocked[id]; ok {
		return ErrWriteLocked
	}
	if c.readLocked[id] > 0 {
		return ErrReadLocked
	}

	c.writeLocked[id] = struct{}{}
	return nil
}

func (c *Cache) writeUnlock(id build.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.writeLocked, id)
}

func (c *Cache) Range(artifactFn func(artifact build.ID) error) error {
	shards, err := ioutil.ReadDir(c.cacheDir)
	if err != nil {
		return err
	}

	for _, shard := range shards {
		dirs, err := ioutil.ReadDir(filepath.Join(c.cacheDir, shard.Name()))
		if err != nil {
			return err
		}

		for _, d := range dirs {
			var id build.ID
			if err := id.UnmarshalText([]byte(d.Name())); err != nil {
				return fmt.Errorf("invalid artifact name: %w", err)
			}

			if err := artifactFn(id); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Cache) Remove(artifact build.ID) error {
	if err := c.writeLock(artifact); err != nil {
		return err
	}
	defer c.writeUnlock(artifact)

	return os.RemoveAll(filepath.Join(c.cacheDir, artifact.Path()))
}

func (c *Cache) Create(artifact build.ID) (path string, commit, abort func() error, err error) {
	if err = c.writeLock(artifact); err != nil {
		return
	}

	path = filepath.Join(c.tmpDir, artifact.String())
	if err = os.MkdirAll(path, 0777); err != nil {
		c.writeUnlock(artifact)
		return
	}

	abort = func() error {
		defer c.writeUnlock(artifact)
		return os.RemoveAll(path)
	}

	commit = func() error {
		defer c.writeUnlock(artifact)
		return os.Rename(path, filepath.Join(c.cacheDir, artifact.Path()))
	}

	return
}

func (c *Cache) Get(artifact build.ID) (path string, unlock func(), err error) {
	if err = c.readLock(artifact); err != nil {
		return
	}

	path = filepath.Join(c.cacheDir, artifact.Path())
	if _, err = os.Stat(path); err != nil {
		c.readUnlock(artifact)

		if os.IsNotExist(err) {
			err = ErrNotFound
		}
		return
	}

	unlock = func() {
		c.readUnlock(artifact)
	}
	return
}

func NewHandler(c *Cache) http.Handler {
	panic("implement me")
}
