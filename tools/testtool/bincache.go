package testtool

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

const BinariesEnv = "TESTTOOL_BINARIES"

type BinCache interface {
	// GetBinary returns filesystem path to the compiled binary corresponding to given import path.
	GetBinary(importPath string) (string, error)
}

type CloseFunc func()

func NewBinCache() (BinCache, CloseFunc) {
	if _, ok := os.LookupEnv(BinariesEnv); ok {
		return newCIBuildCache(), func() {}
	}

	dir, err := ioutil.TempDir("", "bincache-")
	if err != nil {
		log.Fatalf("unable to create temp dir: %s", err)
	}

	return newLocalBinCache(dir), func() { _ = os.RemoveAll(dir) }
}

// localBinCache is a BinCache implementation that compiles queried binaries lazily.
type localBinCache struct {
	// dir is a directory that stores compiled binaries.
	dir string

	binaries sync.Map
}

// newLocalBinCache creates localBinCache that uses given directory to store binaries.
func newLocalBinCache(dir string) *localBinCache {
	return &localBinCache{dir: dir}
}

func (c *localBinCache) GetBinary(importPath string) (string, error) {
	v, ok := c.binaries.Load(importPath)
	if ok {
		return v.(string), nil
	}

	binPath := filepath.Join(c.dir, RandomBinaryName())
	if buildTags == "" {
		runGo("build", "-mod", "readonly", "-o", binPath, importPath)
	} else {
		runGo("build", "-mod", "readonly", "-tags", buildTags, "-o", binPath, importPath)
	}

	c.binaries.Store(importPath, binPath)

	return binPath, nil
}

func runGo(arg ...string) {
	cmd := exec.Command("go", arg...)
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// ciBuildCache is a BinCache implementation that uses precompiled binaries.
type ciBuildCache struct {
	binaries map[string]string
}

// newCIBuildCache creates ciBuildCache that loads locations of precompiled binaries from env variable.
func newCIBuildCache() *ciBuildCache {
	binariesJSON, ok := os.LookupEnv(BinariesEnv)
	if !ok {
		log.Fatalf("%s env variable not set", BinariesEnv)
	}

	binaries := make(map[string]string)
	if err := json.Unmarshal([]byte(binariesJSON), &binaries); err != nil {
		log.Fatalf("unexpected %s format: %s", binaries, err)
	}

	return &ciBuildCache{binaries: binaries}
}

func (c *ciBuildCache) GetBinary(importPath string) (string, error) {
	binary, ok := c.binaries[importPath]
	if !ok {
		return "", fmt.Errorf("%s not found", importPath)
	}
	return binary, nil
}

func RandomBinaryName() string {
	name := RandomName()
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

func RandomName() string {
	var raw [8]byte
	_, _ = rand.Read(raw[:])
	name := hex.EncodeToString(raw[:])
	return name
}
