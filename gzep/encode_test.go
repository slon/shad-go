package gzep_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/gzep"
	"gitlab.com/slon/shad-go/tools/testtool"
)

func BenchmarkEncode(b *testing.B) {
	data := []byte(testtool.RandomName() +
		"New function should generally only return pointer types, " +
		"since a pointer can be put into the return interface " +
		"value without an allocation.",
	)

	
	b.ReportAllocs()
	for b.Loop() {
		require.NoError(b, gzep.Encode(data, io.Discard))
	}
}

func TestEncode_RoundTrip(t *testing.T) {
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
			err := gzep.Encode([]byte(tc.in), buf)
			require.NoError(t, err)

			out, err := decode(buf)
			require.NoError(t, err, tc.in)
			require.Equal(t, tc.in, string(out))
		})

	}
}

func TestEncode_Stress(t *testing.T) {
	wg := &sync.WaitGroup{}

	n := 100
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_ = gzep.Encode([]byte("stonks"), io.Discard)
		}()
	}

	wg.Wait()
}

func TestEncode_Compression(t *testing.T) {
	buf := new(bytes.Buffer)
	err := gzep.Encode(bytes.Repeat([]byte{0x1f}, 1000), buf)
	require.NoError(t, err)
	require.Less(t, buf.Len(), 1000)
}

func decode(r io.Reader) ([]byte, error) {
	rr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rr.Close() }()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, rr); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
