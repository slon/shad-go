//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #include <openssl/blowfish.h>
import "C"

type Blowfish struct {
}

func New(key []byte) *Blowfish {
	panic("implement me")
}
