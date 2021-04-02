package dynamic_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/dynamic"
	"github.com/stretchr/testify/require"
)

type Obj struct {
	V string `json:"v"`
}

func TestJSON(t *testing.T) {
	assert := require.New(t)
	data, err := json.Marshal(nil)
	assert.NoError(err)
	assert.Equal("null", string(data))
	jd := dynamic.JSON(data)

	assert.True(jd.IsNull())
	assert.False(jd.IsEmptyObject())
	assert.False(jd.IsObject())
	assert.False(jd.IsEmptyArray())
	assert.False(jd.IsArray())

	data, err = json.Marshal([]string{"1,2,3"})
	assert.NoError(err)
	jd = dynamic.JSON(data)
	assert.True(jd.IsArray())
	assert.False(jd.IsNull())
	assert.False(jd.IsEmptyObject())
	assert.False(jd.IsObject())
	assert.False(jd.IsEmptyArray())

	data, err = json.Marshal([]string{})
	assert.NoError(err)
	jd = dynamic.JSON(data)
	assert.True(jd.IsArray())
	assert.False(jd.IsNull())
	assert.False(jd.IsEmptyObject())
	assert.False(jd.IsObject())
	assert.True(jd.IsEmptyArray())

	data, err = json.Marshal(map[string]string{})
	assert.NoError(err)
	jd = dynamic.JSON(data)
	assert.False(jd.IsArray())
	assert.False(jd.IsNull())
	assert.True(jd.IsEmptyObject())
	assert.True(jd.IsObject())
	assert.False(jd.IsEmptyArray())

	data, err = json.Marshal(map[string]string{"key": "val"})
	assert.NoError(err)
	jd = dynamic.JSON(data)
	assert.False(jd.IsArray())
	assert.False(jd.IsNull())
	assert.False(jd.IsEmptyObject())
	assert.True(jd.IsObject())
	assert.False(jd.IsEmptyArray())

}

func TestJSONOBject(t *testing.T) {
	assert := require.New(t)
	o := Obj{V: "value"}

	od, err := json.Marshal(o)
	assert.NoError(err)
	obj := dynamic.JSONObject{
		"key": od,
	}
	objData, err := json.Marshal(obj)
	assert.NoError(err)
	assert.Equal(`{"key":{"v":"value"}}`, string(objData))
}
