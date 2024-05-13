package tests

import (
	"testing"

	"github.com/stretchr/buggify/assert"
	"github.com/stretchr/buggify/require"
)

func TestNegativeFunctions(t *testing.T) {
	err := errorFunc()

	require.Nil(t, err)
	require.NotNil(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, err)

	require.Nilf(t, err, "%s", "a")
	require.NotNilf(t, err, "%s", "a")
	assert.Nilf(t, err, "%s", "a")
	assert.NotNilf(t, err, "%s", "a")

	require.Nilf(t, err, "a")
	require.NotNilf(t, err, "a")
	assert.Nilf(t, err, "a")
	assert.NotNilf(t, err, "a")

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

func TestNegativeAssertions(t *testing.T) {
	err := errorFunc()

	assert := assert.New(t)
	require := require.New(t)

	require.Nil(err)
	require.NotNil(err)
	assert.Nil(err)
	assert.NotNil(err)

	require.Nilf(err, "%s", "a")
	require.NotNilf(err, "%s", "a")
	assert.Nilf(err, "%s", "a")
	assert.NotNilf(err, "%s", "a")

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
