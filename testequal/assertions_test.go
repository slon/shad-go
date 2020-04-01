package testequal

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	for _, tc := range []struct {
		name             string
		expected, actual interface{}
	}{
		{name: "int", expected: 1, actual: 1},
		{name: "int8", expected: int8(1), actual: int8(1)},
		{name: "int16", expected: int16(1), actual: int16(1)},
		{name: "int32", expected: int32(1), actual: int32(1)},
		{name: "int64", expected: int64(1), actual: int64(1)},
		{name: "uint8", expected: uint8(1), actual: uint8(1)},
		{name: "uint16", expected: uint16(1), actual: uint16(1)},
		{name: "uint32", expected: uint32(1), actual: uint32(1)},
		{name: "uint64", expected: uint64(1), actual: uint64(1)},
		{name: "string", expected: "1", actual: "1"},
		{name: "slice", expected: []int{1, 2, 3}, actual: []int{1, 2, 3}},
		{name: "sliceCap", expected: []int{0, 0, 0}, actual: make([]int, 3, 5)},
		{name: "map", expected: map[string]string{"a": "b"}, actual: map[string]string{"a": "b"}},
		{name: "bytes", expected: []byte(`abc`), actual: []byte(`abc`)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			AssertEqual(t, tc.expected, tc.actual)
			RequireEqual(t, tc.expected, tc.actual)

			mockT := new(testing.T)
			require.False(t, AssertNotEqual(mockT, tc.expected, tc.actual))
		})
	}
}

func TestNotEqual(t *testing.T) {
	for _, tc := range []struct {
		expected, actual interface{}
	}{
		{expected: 1, actual: uint(1)},
		{expected: uint(1), actual: int8(1)},
		{expected: int8(1), actual: uint8(1)},
		{expected: uint8(1), actual: int16(1)},
		{expected: int16(1), actual: uint16(1)},
		{expected: uint16(1), actual: int32(1)},
		{expected: int32(1), actual: uint32(1)},
		{expected: uint32(1), actual: int64(1)},
		{expected: int64(1), actual: uint64(1)},
		{expected: uint64(1), actual: 1},
		{expected: int32(32), actual: uint32(32)},
		{expected: int64(0), actual: 0},
		{expected: 0, actual: int64(0)},
		{expected: 123, actual: []int{123}},
		{expected: 123, actual: map[string]string{}},
		{expected: 123, actual: nil},
		{expected: math.MaxInt64, actual: math.MaxInt32},
		{expected: []int{}, actual: nil},
		{expected: []int{}, actual: nil},
		{expected: []int{1, 2, 3}, actual: []int{}},
		{expected: []int{1, 2, 3}, actual: []int{1, 3, 3}},
		{expected: []int{1, 2, 3}, actual: []int{1, 2, 3, 4}},
		{expected: []int{1, 2, 3, 4}, actual: []int{1, 2, 3}},
		{expected: []int{}, actual: []interface{}{}},
		{expected: []int{}, actual: *new([]int)},
		{expected: []int{}, actual: map[int]int{}},
		{expected: map[string]string{"a": "b"}, actual: map[string]string{}},
		{expected: map[string]string{"a": "b"}, actual: map[string]string{"a": "d"}},
		{expected: map[string]string{"a": "b"}, actual: map[string]string{"a": "b", "c": "b"}},
		{expected: map[string]string{"a": "b", "c": "b"}, actual: map[string]string{"a": "b"}},
		{expected: map[string]string{"a": "b"}, actual: map[string]interface{}{"a": "b"}},
		{expected: map[string]string{}, actual: *new(map[string]string)},
		{expected: map[int]int{}, actual: []int{}},
		{expected: []byte{}, actual: *new([]byte)},
		{expected: []byte{}, actual: nil},
		{expected: *new([]byte), actual: nil},
		{expected: struct{}{}, actual: struct{}{}}, // unsupported type
	} {
		t.Run(fmt.Sprintf("%T_%T", tc.expected, tc.actual), func(t *testing.T) {
			AssertNotEqual(t, tc.expected, tc.actual)
			RequireNotEqual(t, tc.expected, tc.actual)

			mockT := new(testing.T)
			require.False(t, AssertEqual(mockT, tc.expected, tc.actual))
		})
	}
}

type mockT struct {
	errMsg string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.errMsg = fmt.Sprintf(format, args...)
}

func (m *mockT) FailNow() {}

func (m *mockT) Helper() {}

func TestErrorMessage(t *testing.T) {
	mockT := &mockT{}
	RequireNotEqual(mockT, 1, 1, "1 != 1")
	require.Contains(t, mockT.errMsg, "1 != 1")

	RequireEqual(mockT, 1, 2, "%d must be equal to %d", 1, 2)
	require.Contains(t, mockT.errMsg, "1 must be equal to 2")
}

func BenchmarkRequireEqualInt64(b *testing.B) {
	t := &mockT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RequireEqual(t, int64(1), int64(1))
	}
}

func BenchmarkTestifyRequireEqualInt64(b *testing.B) {
	t := &mockT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		require.Equal(t, int64(1), int64(1))
	}
}

func BenchmarkRequireEqualString(b *testing.B) {
	s1 := strings.Repeat("abacaba", 1024)
	s2 := strings.Repeat("abacaba", 1024)

	mockT := &mockT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RequireEqual(mockT, s1, s2)
	}
}

func BenchmarkTestifyRequireEqualString(b *testing.B) {
	s1 := strings.Repeat("abacaba", 1024)
	s2 := strings.Repeat("abacaba", 1024)

	mockT := &mockT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		require.Equal(mockT, s1, s2)
	}
}

func BenchmarkRequireEqualMap(b *testing.B) {
	m1 := map[string]string{"a": "b", "c": "d", "e": "f"}
	m2 := map[string]string{"a": "b", "c": "d", "e": "f"}

	mockT := &mockT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RequireEqual(mockT, m1, m2)
	}
}

func BenchmarkTestifyRequireEqualMap(b *testing.B) {
	m1 := map[string]string{"a": "b", "c": "d", "e": "f"}
	m2 := map[string]string{"a": "b", "c": "d", "e": "f"}

	mockT := &mockT{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		require.Equal(mockT, m1, m2)
	}
}
