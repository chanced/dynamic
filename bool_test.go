package dynamic_test

import (
	"testing"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	assert := require.New(t)
	b := dynamic.NewBool("true")
	v, ok := b.Bool()
	assert.True(ok, "should have a bool value")
	assert.True(v, "should be set to true")

	b = dynamic.NewBool("false")
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.False(v, "should be set to false")

	b = dynamic.NewBool(true)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.True(v, "should be set to true")

	b = dynamic.NewBool(false)
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.False(v, "should be set to false")

	b = dynamic.NewBool("1")
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.True(v, "should be set to true")

	b = dynamic.NewBool("0")
	v, ok = b.Bool()
	assert.True(ok, "should have a bool value")
	assert.False(v, "should be set to false")

}
