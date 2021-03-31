package dynamic

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	trueBytes  = []byte("true")
	falseBytes = []byte("false")
)

type JSON []byte

func (d JSON) Len() int {
	return len(d)
}
func (d JSON) MarshalJSON() ([]byte, error) {
	if d == nil || len(d) == 0 {
		return Null, nil
	}

	return d, nil

}
func (d *JSON) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("dynamic.RawMessage: UnmarshalJSON on nil pointer")
	}
	*d = append((*d)[0:0], data...)
	return nil
}
func (d JSON) IsObject() bool {
	if len(d) < 2 {
		return false
	}
	return d[0] == '{' && d[len(d)-1] == '}'
}

func (d JSON) IsEmptyObject() bool {
	return d.IsObject() && len(d) == 2
}

func (d JSON) IsEmptyArray() bool {
	return d.IsArray() && len(d) == 2
}

// IsArray reports whether the data is a json array. It does not check whether
// the json is malformed.
func (d JSON) IsArray() bool {
	if len(d) < 2 {
		return false
	}
	return d[0] == '[' && d[len(d)-1] == ']'
}

func (d JSON) IsNull() bool {
	return bytes.Equal(d, Null)
}

// IsBool reports true if data appears to be a json boolean value. It is
// possible that it will report false positives of malformed json.
//
// IsBool does not parse strings
func (d JSON) IsBool() bool {
	if len(d) < 4 {
		return false
	}
	return (d[0] == 't' && len(d) == 4) || (d[0] == 'f' && len(d) == 5)
}

// IsTrue reports true if data appears to be a json boolean value of true. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsTrue does not parse strings
func (d JSON) IsTrue() bool {
	return d.IsBool() && d[0] == 't'
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func (d JSON) IsFalse() bool {
	return bytes.Equal(d, falseBytes)
}

func (d JSON) Equal(data []byte) bool {
	return bytes.Equal(d, data)
}

// ContainsEscapeRune reports whether the string value of d contains "\"
// It returns false if d is not a quoted string.
func (d JSON) ContainsEscapeRune() bool {
	for i := 0; i < len(d); i++ {
		if d[i] == '\\' {
			return true
		}
	}
	return false
}

// UnquotedString trims double quotes from the bytes. It does not parse for
// escaped characters
func (d JSON) UnquotedString() string {
	if len(d) < 2 {
		return string(d)
	}

	if d[0] == '"' && d[len(d)-1] == '"' {
		return string(d[1 : len(d)-1])
	}
	return string(d)
}

func (d JSON) IsNumber() bool {
	if len(d) == 0 {
		return false
	}
	switch d[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return true
	default:
		return false
	}
}
func (d JSON) IsString() bool {
	if len(d) == 0 {
		return false
	}
	return d[0] == '"'
}

type JSONObject map[string]JSON

func (obj JSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]JSON(obj))
}

func (obj *JSONObject) UnmarshalJSON(data []byte) error {
	var m map[string]JSON
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	*obj = m
	return nil
}
