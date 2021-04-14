package illegal_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"gitlab.com/slon/shad-go/illegal"
)

func TestStringFromBytes(t *testing.T) {
	b := []byte{'a', 'b', 'c'}
	s := illegal.StringFromBytes(b)

	assert.Equal(t, "abc", s)
	assert.Equal(t, *(*uintptr)(unsafe.Pointer(&b)), *(*uintptr)(unsafe.Pointer(&s)))
}
