package illegal_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"gitlab.com/slon/shad-go/illegal"
)

func TestStringFromBytes(t *testing.T) {
	var tests = []struct {
		name  string
		input []byte
	}{
		{"NotEmpty", []byte{'a', 'b', 'c'}},
		{"Empty", []byte{}},
		{"Nil", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := test.input
			s := illegal.StringFromBytes(b)

			bptr := *(*uintptr)(unsafe.Pointer(&b))
			sptr := *(*uintptr)(unsafe.Pointer(&s))

			assert.Equal(t, string(b), s)
			assert.Equal(t, bptr, sptr, "string ptr [%v] != []byte ptr [%v]\n", sptr, bptr)
		})
	}
}
