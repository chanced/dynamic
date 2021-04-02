package dynamic_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

func TestNumber(t *testing.T) {
	assert := require.New(t)
	n, err := dynamic.NewNumber(34.34)
	assert.NoError(err)

	f, ok := n.Float64()
	assert.True(ok, "should be able to retrieve a float value")
	assert.Equal(f, float64(34.34))
	_, ok = n.Int64()
	assert.False(ok, "should not be able to get an int64 from a float")
	_, ok = n.Uint64()
	assert.False(ok, "should not be able to get an uint64 from a float")

	data, err := json.Marshal(f)
	assert.NoError(err)
	assert.Equal("34.34", string(data))
	var nn dynamic.Number
	err = json.Unmarshal([]byte("34.34"), &nn)
	assert.NoError(err)
	f, ok = nn.Float64()
	assert.True(ok, "should be able to retrieve a float value")
	assert.Equal(f, float64(34.34))

	err = json.Unmarshal([]byte("0.001"), &nn)
	assert.NoError(err)
	f, ok = nn.Float64()
	assert.True(ok, "should be able to retrieve a float value")
	assert.Equal(f, float64(0.001))
	assert.True(nn.HasValue())

	ni, err := dynamic.NewNumber(34)
	assert.NoError(err)
	i, ok := ni.Int64()
	assert.True(ok)
	assert.Equal(int64(34), i)
	str := ni.String()
	assert.NotEmpty(str)
}
