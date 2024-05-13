package tests

import (
	"testing"

	xassert "github.com/stretchr/testify/assert"
	xrequire "github.com/stretchr/testify/require"
)

func TestXFunctions(t *testing.T) {
	err := errorFunc()

	xrequire.Nil(t, err)    // want `use require.NoError instead of comparing error to nil`
	xrequire.NotNil(t, err) // want `use require.Error instead of comparing error to nil`
	xassert.Nil(t, err)     // want `use assert.NoError instead of comparing error to nil`
	xassert.NotNil(t, err)  // want `use assert.Error instead of comparing error to nil`

	xrequire.Nilf(t, err, "%s", "a")    // want `use require.NoErrorf instead of comparing error to nil`
	xrequire.NotNilf(t, err, "%s", "a") // want `use require.Errorf instead of comparing error to nil`
	xassert.Nilf(t, err, "%s", "a")     // want `use assert.NoErrorf instead of comparing error to nil`
	xassert.NotNilf(t, err, "%s", "a")  // want `use assert.Errorf instead of comparing error to nil`

	xrequire.Nilf(t, err, "a")    // want `use require.NoErrorf instead of comparing error to nil`
	xrequire.NotNilf(t, err, "a") // want `use require.Errorf instead of comparing error to nil`
	xassert.Nilf(t, err, "a")     // want `use assert.NoErrorf instead of comparing error to nil`
	xassert.NotNilf(t, err, "a")  // want `use assert.Errorf instead of comparing error to nil`

	p := new(int)

	xrequire.Nil(t, p)
	xrequire.NotNil(t, p)
	xassert.Nil(t, p)
	xassert.NotNil(t, p)

	xrequire.Nilf(t, p, "%s", "a")
	xrequire.NotNilf(t, p, "%s", "a")
	xassert.Nilf(t, p, "%s", "a")
	xassert.NotNilf(t, p, "%s", "a")
}

func TestXAssertions(t *testing.T) {
	err := errorFunc()

	xassert := xassert.New(t)
	xrequire := xrequire.New(t)

	xrequire.Nil(err)    // want `use require.NoError instead of comparing error to nil`
	xrequire.NotNil(err) // want `use require.Error instead of comparing error to nil`
	xassert.Nil(err)     // want `use assert.NoError instead of comparing error to nil`
	xassert.NotNil(err)  // want `use assert.Error instead of comparing error to nil`

	xrequire.Nilf(err, "%s", "a")    // want `use require.NoErrorf instead of comparing error to nil`
	xrequire.NotNilf(err, "%s", "a") // want `use require.Errorf instead of comparing error to nil`
	xassert.Nilf(err, "%s", "a")     // want `use assert.NoErrorf instead of comparing error to nil`
	xassert.NotNilf(err, "%s", "a")  // want `use assert.Errorf instead of comparing error to nil`

	p := new(int)

	xrequire.Nil(p)
	xrequire.NotNil(p)
	xassert.Nil(p)
	xassert.NotNil(p)

	xrequire.Nilf(p, "%s", "a")
	xrequire.NotNilf(p, "%s", "a")
	xassert.Nilf(p, "%s", "a")
	xassert.NotNilf(p, "%s", "a")
}
