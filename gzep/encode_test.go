package gzep

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/tools/testtool"
)

func BenchmarkEncode(b *testing.B) {
	data := []byte(
		"New function should generally only return pointer types, " +
			"since a pointer can be put into the return interface " +
			"value without an allocation. " +
			testtool.RandomName(),
	)

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		require.NoError(b, Encode(data, io.Discard))
	}
}

func TestRoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		in   string
	}{
		{name: "empty", in: ""},
		{name: "simple", in: "A long time ago in a galaxy far, far away..."},
		{name: "random", in: testtool.RandomName()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			require.NoError(t, Encode([]byte(tc.in), buf))

			out, err := decode(buf.Bytes())
			require.NoError(t, err, tc.in)
			require.Equal(t, tc.in, string(out))
		})

	}
}

func decode(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
