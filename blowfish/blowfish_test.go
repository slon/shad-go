package blowfish_test

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/blowfish"
)

var _ cipher.Block = (*blowfish.Blowfish)(nil)

func TestBlowfish(t *testing.T) {
	b := blowfish.New([]byte("kek"))

	for i, testCase := range []struct {
		in  uint64
		enc uint64
		dec uint64
	}{
		{
			in:  0x0,
			enc: 0x03e009b8123919ea,
			dec: 0xc5b3bba65042b0bf,
		},
		{
			in:  0x0123456789abcdef,
			enc: 0x1c7879d650892fe0,
			dec: 0xf714799fdf68637c,
		},
	} {

		var in, out, expected [8]byte

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			binary.BigEndian.PutUint64(in[:], testCase.in)

			b.Encrypt(out[:], in[:])
			binary.BigEndian.PutUint64(expected[:], testCase.enc)
			require.Equal(t, expected, out)

			b.Decrypt(out[:], in[:])
			binary.BigEndian.PutUint64(expected[:], testCase.dec)
			require.Equal(t, expected, out)
		})
	}
}
