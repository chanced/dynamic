package dynamic_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

type Struct struct {
	V1 dynamic.BoolOrString `json:"v1"`
	V2 dynamic.BoolOrString `json:"v2,omitempty"`
}

func TestBoolOrString(t *testing.T) {
	raw1 := []byte(`{
		"v1": true,
		"v2": "truth"
	}`)
	assert := require.New(t)

	var s1 Struct
	err := json.Unmarshal(raw1, &s1)
	assert.NoError(err)
	bt, err := dynamic.NewBoolOrString("true")
	assert.NoError(err)
	assert.Equal(bt, s1.V1)
	tr, err := dynamic.NewBoolOrString("truth")
	assert.NoError(err)
	assert.Equal(tr, s1.V2)

	m1, err := json.Marshal(s1)
	assert.NoError(err)
	assert.Equal(`{"v1":true,"v2":"truth"}`, string(m1))
	raw2 := []byte(`{
		"v1": false,
		"v2": "STRVAL"
	}`)

	var s2 Struct
	err = json.Unmarshal(raw2, &s2)
	assert.NoError(err)
	f, err := dynamic.NewBoolOrString("false")
	assert.NoError(err)
	assert.Equal(f, s2.V1)
	strv, err := dynamic.NewBoolOrString("STRVAL")
	assert.NoError(err)
	assert.Equal(strv, s2.V2)

	m2, err := json.Marshal(s2)
	assert.NoError(err)

	fmt.Println(string(m2))
	assert.Equal(`{"v1":false,"v2":"STRVAL"}`, string(m2))
}
