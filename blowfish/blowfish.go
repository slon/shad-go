//go:build !solution

package blowfish

// #cgo LDFLAGS: -lcrypto
// #include <openssl/blowfish.h>
import "C"

type Blowfish struct {
}

func New(key []byte) *Blowfish {
	panic("implement me")
}
