package externalsort

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func newStringReader(s string) LineReader {
	return NewReader(strings.NewReader(s))
}

func readAll(r LineReader) (lines []string, err error) {
	for {
		l, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if l != "" {
					lines = append(lines, l)
				}
				return lines, nil
			}
			return nil, err
		}
		lines = append(lines, l)
	}
}

func TestLineReader(t *testing.T) {
	type Wrapper func(r io.Reader) io.Reader

	for _, tc := range []struct {
		name     string
		in       string
		wrappers []Wrapper
		expected []string
	}{
		{
			name:     "empty",
			in:       "",
			expected: nil,
		},
		{
			name:     "one-row",
			in:       "abc",
			expected: []string{"abc"},
		},
		{
			name: "linebreak",
			in: `abc

`,
			expected: []string{"abc", ""},
		},
		{
			name: "multiple-rows",
			in: `a

b
b
`,
			expected: []string{"a", "", "b", "b"},
		},
		{
			name:     "large-row",
			in:       strings.Repeat("a", 4097),
			expected: []string{strings.Repeat("a", 4097)},
		},
		{
			name:     "huge-row",
			in:       strings.Repeat("a", 65537),
			expected: []string{strings.Repeat("a", 65537)},
		},
		{
			name:     "half-reader",
			in:       strings.Repeat("a", 1025),
			wrappers: []Wrapper{iotest.HalfReader},
			expected: []string{strings.Repeat("a", 1025)},
		},
		{
			name:     "eof",
			in:       strings.Repeat("a", 1025),
			wrappers: []Wrapper{iotest.DataErrReader},
			expected: []string{strings.Repeat("a", 1025)},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var r io.Reader
			r = strings.NewReader(tc.in)
			for _, w := range tc.wrappers {
				r = w(r)
			}

			lineReader := NewReader(r)

			lines, err := readAll(lineReader)
			require.NoError(t, err)

			require.Len(t, lines, len(tc.expected),
				"expected: %+v, got: %+v", tc.expected, lines)
			require.Equal(t, tc.expected, lines)
		})
	}
}

type brokenReader int

func (r brokenReader) Read(data []byte) (n int, err error) {
	return 0, errors.New("read is broken")
}

type eofReader int

func (r eofReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func TestLineReader_error(t *testing.T) {
	_, err := NewReader(new(brokenReader)).ReadLine()
	require.Error(t, err)
	require.False(t, errors.Is(err, io.EOF))

	r := NewReader(new(eofReader))
	_, err = r.ReadLine()
	require.True(t, errors.Is(err, io.EOF))
}

func TestLineWriter(t *testing.T) {
	for _, tc := range []struct {
		name  string
		lines []string
	}{
		{
			name:  "empty",
			lines: []string{""},
		},
		{
			name:  "simple",
			lines: []string{"a", "b", "c"},
		},
		{
			name:  "large-line",
			lines: []string{strings.Repeat("xx", 2049), "x", "y"},
		},
		{
			name:  "huge-line",
			lines: []string{strings.Repeat("?", 65537), "?", "!"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			lw := NewWriter(w)

			for _, l := range tc.lines {
				require.NoError(t, lw.Write(l))
			}

			require.NoError(t, w.Flush())
			expected := strings.Join(tc.lines, "\n") + "\n"
			require.Equal(t, expected, buf.String())
		})
	}
}
