package example

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestFoo(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := NewMockFoo(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	m.
		EXPECT().
		Bar(gomock.Eq(99)).
		Return(101)

	SUT(m)
}
