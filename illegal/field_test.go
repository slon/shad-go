package illegal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/slon/shad-go/illegal"
	"gitlab.com/slon/shad-go/illegal/internal"
)

func TestIllegalField(t *testing.T) {
	var s internal.Struct

	illegal.SetPrivateField(&s, "a", 10)
	illegal.SetPrivateField(&s, "b", "foo")
	illegal.SetPrivateField(&s, "p", internal.NewPrivateType(42))

	assert.Equal(t, "10 foo 42", s.String())
}

func TestIllegalWrongFieldType(t *testing.T) {
	var s internal.Struct

	assert.Panics(t, func() {
		illegal.SetPrivateField(&s, "a", "1234")
	})
}
