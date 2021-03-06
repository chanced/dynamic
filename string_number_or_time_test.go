package dynamic_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

type SNT struct {
	V *dynamic.StringNumberOrTime `json:"v,omitempty"`
}

func TestStringNumberOrTime(t *testing.T) {
	assert := require.New(t)
	var expectedJSON []byte
	p1 := SNT{}
	p2 := SNT{}
	now, err := time.Parse(dynamic.DefaultTimeLayout(), time.Now().Format(dynamic.DefaultTimeLayout()))
	assert.NoError(err)
	p1.V = &dynamic.StringNumberOrTime{}
	err = p1.V.Set(now)
	assert.NoError(err)
	tv, isTime := p1.V.Time()
	assert.True(isTime)
	assert.Equal(now, tv)

	p1.V = &dynamic.StringNumberOrTime{}
	err = p1.V.Set(now.Format(dynamic.DefaultTimeLayout()))
	assert.NoError(err)
	tv, isTime = p1.V.Time()
	assert.True(isTime)
	assert.Equal(now.Format(dynamic.DefaultTimeLayout()), p1.V.String())

	expectedJSON = []byte(`{"v":"` + now.Format(dynamic.DefaultTimeLayout()) + `"}`)

	b, err := json.Marshal(p1)
	assert.NoError(err)
	assert.Equal(expectedJSON, b)

	err = json.Unmarshal(expectedJSON, &p2)
	assert.NoError(err)

	tm, ok := p2.V.Time()
	assert.True(ok)
	assert.Equal(now, tm)

	assert.NotEmpty(p1.V.String())
	err = p1.V.Set(uint64(0xFFFFFFFFFFFFFFFF))
	assert.NoError(err)
	uv, ok := p1.V.Uint64()
	assert.True(ok)
	assert.Equal(uint64(0xFFFFFFFFFFFFFFFF), uv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":"18446744073709551615"}`)
	assert.Equal(expectedJSON, b)

	err = p1.V.Set(uint64(234))
	assert.NoError(err)
	uv, ok = p1.V.Uint64()
	assert.True(ok)
	assert.Equal(uint64(234), uv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":234}`)
	assert.Equal(expectedJSON, b)

	p1.V.Set(int64(234))
	iv, ok := p1.V.Int64()
	assert.True(ok)
	assert.Equal(int64(234), iv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":234}`)
	assert.Equal(expectedJSON, b)

	_, ok = p1.V.Time()
	assert.False(ok)
	assert.False(p1.V.IsTime())

	p1.V.Set("0xFFFFFFFFFFFFFF00")
	uv, ok = p1.V.Uint64()
	assert.True(ok)
	bigu := uint64(0xFFFFFFFFFFFFFF00)
	assert.Equal(bigu, uv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":"18446744073709551360"}`)
	assert.Equal(expectedJSON, b)
	err = json.Unmarshal(expectedJSON, &p2)
	assert.NoError(err)
	ui, ok := p2.V.Uint64()
	assert.True(ok)
	assert.Equal(bigu, ui)

	_, ok = p1.V.Time()
	assert.False(ok)
	assert.False(p1.V.IsTime())

	i := int64(-1239122)
	p1.V.Set(i)
	iv, ok = p1.V.Int64()
	assert.True(ok)
	assert.Equal(i, iv)
	uv, ok = p1.V.Uint64()
	assert.False(ok)
	assert.Zero(uv)

	p1.V.Set(958.34)
	f, ok := p1.V.Float64()
	assert.True(ok)
	assert.Equal(958.34, f)

	p1.V.Set(float64(-3942.2))
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":-3942.2}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)

	assert.NoError(err)
	f2, ok := p2.V.Float64()
	assert.True(ok)
	assert.Equal(f2, -3942.2)

	err = p1.V.Set(float64(-992345.1233))
	assert.NoError(err)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":-992345.1233}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)

	assert.NoError(err)
	f2, ok = p2.V.Float64()
	assert.True(ok)
	assert.Equal(f2, -992345.1233)

	err = p1.V.Set(float64(-993456789012345.1))
	assert.NoError(err)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":-993456789012345.1}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)

	assert.NoError(err)
	f2, ok = p2.V.Float64()
	assert.True(ok)
	assert.Equal(f2, -993456789012345.1)

}
