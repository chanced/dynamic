package dynamic_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

type SNBT struct {
	V *dynamic.StringNumberBoolOrTime `json:"v,omitempty"`
}

func TestStringNumberBoolOrTime(t *testing.T) {
	assert := require.New(t)
	var expectedJSON []byte
	p1 := SNBT{}
	p2 := SNBT{}
	now, err := time.Parse(dynamic.DefaultTimeLayout(), time.Now().Format(dynamic.DefaultTimeLayout()))
	assert.NoError(err)
	p1.V = dynamic.NewStringNumberBoolOrTime()
	p1.V.Set(now)
	tv, isTime := p1.V.Time()
	assert.True(isTime)
	assert.Equal(now, tv)

	p1.V = dynamic.NewStringNumberBoolOrTime()
	p1.V.Set(now.Format(dynamic.DefaultTimeLayout()))
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
	p1.V.Set(uint64(0xFFFFFFFFFFFFFFFF))
	uv, ok := p1.V.Uint()
	assert.True(ok)
	assert.Equal(uint64(0xFFFFFFFFFFFFFFFF), uv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":"18446744073709551615"}`)
	assert.Equal(expectedJSON, b)

	p1.V.Set(uint64(234))
	uv, ok = p1.V.Uint()
	assert.True(ok)
	assert.Equal(uint64(234), uv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":234}`)
	assert.Equal(expectedJSON, b)

	p1.V.Set(int64(234))
	iv, ok := p1.V.Int()
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
	uv, ok = p1.V.Uint()
	assert.True(ok)
	bigu := uint64(0xFFFFFFFFFFFFFF00)
	assert.Equal(bigu, uv)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":"18446744073709551360"}`)
	assert.Equal(expectedJSON, b)
	err = json.Unmarshal(expectedJSON, &p2)
	assert.NoError(err)
	ui, ok := p2.V.Uint()
	assert.True(ok)
	assert.Equal(bigu, ui)

	_, ok = p1.V.Time()
	assert.False(ok)
	assert.False(p1.V.IsTime())

	i := int64(-1239122)
	p1.V.Set(i)
	iv, ok = p1.V.Int()
	assert.True(ok)
	assert.Equal(i, iv)
	uv, ok = p1.V.Uint()
	assert.False(ok)
	assert.Zero(uv)

	p1.V.Set(958.34)
	f, ok := p1.V.Float()
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
	f2, ok := p2.V.Float()
	assert.True(ok)
	assert.Equal(f2, -3942.2)

	p1.V.Set(float64(-992345.1233))
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":-992345.1233}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)

	assert.NoError(err)
	f2, ok = p2.V.Float()
	assert.True(ok)
	assert.Equal(f2, -992345.1233)

	p1.V.Set(float64(-993456789012345.1))
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":-993456789012345.1}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)

	assert.NoError(err)
	f2, ok = p2.V.Float()
	assert.True(ok)
	assert.Equal(f2, -993456789012345.1)

	err = p1.V.Set(true)
	assert.NoError(err)

	b, err = json.Marshal(p1)
	assert.NoError(err)

	expectedJSON = []byte(`{"v":true}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)

	assert.NoError(err)
	bv, ok := p2.V.Bool()
	assert.True(ok)
	assert.True(bv)

	err = p1.V.Set("true")
	assert.NoError(err)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":"true"}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)
	assert.NoError(err)
	bv, ok = p2.V.Bool()
	assert.True(ok)
	assert.True(bv)

	err = p1.V.Set("false")
	assert.NoError(err)
	b, err = json.Marshal(p1)
	assert.NoError(err)
	expectedJSON = []byte(`{"v":"false"}`)
	assert.Equal(expectedJSON, b)
	assert.NotEmpty(expectedJSON)
	err = json.Unmarshal([]byte(expectedJSON), &p2)
	assert.NoError(err)
	bv, ok = p2.V.Bool()
	assert.True(ok)
	assert.False(bv)

}
