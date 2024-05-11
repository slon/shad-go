//go:build !change

package psycher

import (
	"crypto/cipher"
	"crypto/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type slowCipher struct {
	keys [][]byte
}

var _ cipher.Block = (*slowCipher)(nil)

func (c slowCipher) BlockSize() int {
	return 16
}

func (c slowCipher) Encrypt(dst, src []byte) {
	var idx int
	for srcI, b := range src {
		for ; b != 0; b /= 2 {
			if b%2 == 1 {
				for i := range dst {
					dst[i] ^= c.keys[idx][i]
				}
			}
			idx++
		}
		idx = (srcI + 1) * 8
	}
}

func (c slowCipher) Decrypt(_, _ []byte) {
	panic("something nerdy")
}

func genKeys() [][]byte {
	var keys [][]byte
	for i := range 128 {
		keys = append(keys, make([]byte, 16))
		_, _ = rand.Read(keys[i])
	}
	return keys
}

func TestCorrectness(t *testing.T) {
	keys := genKeys()

	solCipher, correctCipher := New(keys), slowCipher{keys: keys}

	for range 1000 {
		src := make([]byte, 16)
		_, _ = rand.Read(src)

		srcCopy := slices.Clone(src)
		dst1, dst2 := make([]byte, 16), make([]byte, 16)

		correctCipher.Encrypt(dst1, src)
		solCipher.Encrypt(dst2, srcCopy)

		require.Equal(t, src, srcCopy, "shouldn't modify src")
		require.Equal(t, dst1, dst2, "wrong Encrypt output")
	}
}

func benchCipher(b *testing.B, constructor func([][]byte) cipher.Block) {
	keys := genKeys()
	c := constructor(keys)

	const blocksLen = 1_000_000
	var blocks [][]byte
	for range blocksLen {
		src := make([]byte, 16)
		_, _ = rand.Read(src)

		blocks = append(blocks, src)
	}
	dst := make([]byte, 16)

	b.ReportAllocs()
	b.SetBytes(16)
	b.ResetTimer()

	for i := range b.N {
		c.Encrypt(dst, blocks[i%blocksLen])
	}
}

func BenchmarkSlowCipher(b *testing.B) {
	benchCipher(b, func(keys [][]byte) cipher.Block {
		c := slowCipher{keys: keys}
		return c
	})
}

func BenchmarkCipher(b *testing.B) {
	benchCipher(b, func(keys [][]byte) cipher.Block {
		c := New(keys)
		return c
	})
}
