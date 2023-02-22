//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"

type Blowfish struct {
}

func New(key []byte) *Blowfish {
	panic("implement me")
}
