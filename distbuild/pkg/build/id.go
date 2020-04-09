package build

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding"
	"encoding/hex"
	"fmt"
	"path/filepath"
)

type ID [sha1.Size]byte

var (
	_ = encoding.TextMarshaler(ID{})
	_ = encoding.TextUnmarshaler(&ID{})
)

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

func (id ID) Path() string {
	return filepath.Join(hex.EncodeToString(id[:1]), hex.EncodeToString(id[:]))
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(id[:])), nil
}

func (id *ID) UnmarshalText(b []byte) error {
	raw, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}

	if len(raw) != len(id) {
		return fmt.Errorf("invalid id size: %q", b)
	}

	copy(id[:], raw)
	return nil
}

func NewID() ID {
	var id ID
	_, err := rand.Read(id[:])
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: %v", err))
	}
	return id
}
