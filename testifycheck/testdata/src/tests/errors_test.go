package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func errorFunc() error {
	panic("implement me")
}

func TestFunctions(t *testing.T) {
	err := errorFunc()

	require.Nil(t, err)    // want `use require.NoError instead of comparing error to nil`
	require.NotNil(t, err) // want `use require.Error instead of comparing error to nil`
	assert.Nil(t, err)     // want `use assert.NoError instead of comparing error to nil`
	assert.NotNil(t, err)  // want `use assert.Error instead of comparing error to nil`

	require.Nilf(t, err, "%s", "a") // want `use require.NoErrorf instead of comparing error to nil`
	require.NotNilf(t, err, "%s", "a") // want `use require.Errorf instead of comparing error to nil`
	assert.Nilf(t, err, "%s", "a")     // want `use assert.NoErrorf instead of comparing error to nil`
	assert.NotNilf(t, err, "%s", "a")  // want `use assert.Errorf instead of comparing error to nil`

	require.Nilf(t, err, "a")    // want `use require.NoErrorf instead of comparing error to nil`
	require.NotNilf(t, err, "a") // want `use require.Errorf instead of comparing error to nil`
	assert.Nilf(t, err, "a")     // want `use assert.NoErrorf instead of comparing error to nil`
	assert.NotNilf(t, err, "a")  // want `use assert.Errorf instead of comparing error to nil`

	p := new(int)

	require.Nil(t, p)
	require.NotNil(t, p)
	assert.Nil(t, p)
	assert.NotNil(t, p)

	require.Nilf(t, p, "%s", "a")
	require.NotNilf(t, p, "%s", "a")
	assert.Nilf(t, p, "%s", "a")
	assert.NotNilf(t, p, "%s", "a")
}

func TestAssertions(t *testing.T) {
	err := errorFunc()

	assert := assert.New(t)
	require := require.New(t)

	require.Nil(err)    // want `use require.NoError instead of comparing error to nil`
	require.NotNil(err) // want `use require.Error instead of comparing error to nil`
	assert.Nil(err)     // want `use assert.NoError instead of comparing error to nil`
	assert.NotNil(err)  // want `use assert.Error instead of comparing error to nil`

	require.Nilf(err, "%s", "a") // want `use require.NoErrorf instead of comparing error to nil`
	require.NotNilf(err, "%s", "a") // want `use require.Errorf instead of comparing error to nil`
	assert.Nilf(err, "%s", "a")     // want `use assert.NoErrorf instead of comparing error to nil`
	assert.NotNilf(err, "%s", "a")  // want `use assert.Errorf instead of comparing error to nil`

	p := new(int)

	require.Nil(p)
	require.NotNil(p)
	assert.Nil(p)
	assert.NotNil(p)

	require.Nilf(p, "%s", "a")
	require.NotNilf(p, "%s", "a")
	assert.Nilf(p, "%s", "a")
	assert.NotNilf(p, "%s", "a")
}
