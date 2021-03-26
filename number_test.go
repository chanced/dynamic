package dynamic_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

func TestNumber(t *testing.T) {
	assert := require.New(t)
	n := dynamic.NewNumber(34.34)

	f, ok := n.Float()
	assert.True(ok, "should be able to retrieve a float value")
	assert.Equal(f, float64(34.34))
	_, ok = n.Int()
	assert.False(ok, "should not be able to get an int64 from a float")
	_, ok = n.Uint()
	assert.False(ok, "should not be able to get an uint64 from a float")

	data, err := json.Marshal(f)
	assert.NoError(err)
	assert.Equal("34.34", string(data))
	var nn dynamic.Number
	err = json.Unmarshal([]byte("34.34"), &nn)
	assert.NoError(err)
	f, ok = nn.Float()
	assert.True(ok, "should be able to retrieve a float value")
	assert.Equal(f, float64(34.34))
}
