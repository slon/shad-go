//go:build !solution

package psycher

import "crypto/cipher"

// Cipher describes psycher block cipher.
type Cipher struct {
}

var _ cipher.Block = (*Cipher)(nil)

// New creates Cipher.
func New(keys [][]byte) *Cipher {
	panic("implement me")
}

// BlockSize returns the cipher's block size.
func (c Cipher) BlockSize() int {
	panic("implement me")
}

// Encrypt encrypts the first block in src into dst.
func (c Cipher) Encrypt(dst, src []byte) {
	panic("implement me")
}

// Decrypt decrypts the first block in src into dst.
func (c Cipher) Decrypt(dst, src []byte) {
	panic("implement me")
}
