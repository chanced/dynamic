package dynamic_test

import (
	"testing"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	assert := require.New(t)
	b, err := dynamic.NewBool("true")
	assert.NoError(err)
	v, ok := b.Bool()
	assert.True(ok, "should have a bool value")
	assert.True(v, "should be set to true")

	b, err = dynamic.NewBool("false")
	assert.NoError(err)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.False(v, "should be set to false")

	b, err = dynamic.NewBool(true)
	assert.NoError(err)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.True(v, "should be set to true")

	b, err = dynamic.NewBool(false)
	assert.NoError(err)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.False(v, "should be set to false")

	b, err = dynamic.NewBool("1")
	assert.NoError(err)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.True(v, "should be set to true")

	b, err = dynamic.NewBool("0")
	assert.NoError(err)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.False(v, "should be set to false")
}
