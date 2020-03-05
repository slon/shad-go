package otp

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

type zeroSometimesReader struct {
	io.Reader
	i int
}

func (r *zeroSometimesReader) Read(p []byte) (n int, err error) {
	r.i++
	if r.i&1 == 0 {
		return 0, nil
	}

	return r.Reader.Read(p)
}

const testSize = 1234

var (
	plaintext   = make([]byte, testSize)
	randomBytes = make([]byte, testSize)
	ciphertext  = make([]byte, testSize)

	plaintextBackup = make([]byte, testSize)
)

func init() {
	rand.Read(plaintext)
	rand.Read(randomBytes)

	copy(plaintextBackup, plaintext)

	for i := range plaintext {
		ciphertext[i] = plaintext[i] ^ randomBytes[i]
	}
}

func TestReader(t *testing.T) {
	for _, testCase := range []struct {
		name string
		r    io.Reader
		prng io.Reader

		err    error
		result []byte
	}{
		{
			name:   "simple",
			r:      bytes.NewBuffer(plaintext),
			prng:   bytes.NewBuffer(randomBytes),
			result: ciphertext,
		},
		{
			name:   "eof",
			r:      iotest.DataErrReader(bytes.NewBuffer(plaintext)),
			prng:   bytes.NewBuffer(randomBytes),
			result: ciphertext,
		},
		{
			name:   "halfreader",
			r:      iotest.HalfReader(bytes.NewBuffer(plaintext)),
			prng:   bytes.NewBuffer(randomBytes),
			result: ciphertext,
		},
		{
			name:   "zerosometimes",
			r:      &zeroSometimesReader{Reader: iotest.HalfReader(iotest.HalfReader(bytes.NewBuffer(plaintext)))},
			prng:   bytes.NewBuffer(randomBytes),
			result: ciphertext,
		},
		{
			name:   "timeout",
			r:      iotest.TimeoutReader(bytes.NewBuffer(plaintext)),
			prng:   bytes.NewBuffer(randomBytes),
			result: ciphertext[:bytes.MinRead],
			err:    iotest.ErrTimeout,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewReader(testCase.r, testCase.prng)

			buf, err := ioutil.ReadAll(r)
			require.Equal(t, testCase.err, err)
			require.Equal(t, testCase.result, buf)
		})
	}
}

func TestWriterSimple(t *testing.T) {
	out := &bytes.Buffer{}
	prng := bytes.NewBuffer(randomBytes)

	w := NewWriter(out, prng)
	n, err := w.Write(plaintext)

	require.Equalf(t, plaintextBackup, plaintext, "Write must not modify the slice data, even temporarily.")
	require.NoError(t, err)
	require.Equal(t, len(plaintext), n)
	require.Equal(t, out.Bytes(), ciphertext)
}

type errWriter struct {
	buf bytes.Buffer
	n   int
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	if len(p) > w.n {
		p = p[:w.n]
	}

	n = len(p)
	w.n -= n

	if w.n == 0 {
		err = iotest.ErrTimeout
	}

	w.buf.Write(p)
	return
}

func TestWriterError(t *testing.T) {
	out := &errWriter{n: 512}
	prng := bytes.NewBuffer(randomBytes)

	w := NewWriter(out, prng)
	n, err := w.Write(plaintext)

	require.Equalf(t, plaintextBackup, plaintext, "Write must not modify the slice data, even temporarily.")
	require.Equal(t, err, iotest.ErrTimeout)
	require.Equal(t, 512, n)
	require.Equal(t, out.buf.Bytes(), ciphertext[:512])
}
