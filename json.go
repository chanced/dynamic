package dynamic

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

var (
	trueBytes  = []byte("true")
	falseBytes = []byte("false")
)

type JSON []byte

func (raw JSON) MarshalJSON() ([]byte, error) {
	if raw == nil {
		return Null, nil
	}
	return raw, nil

}
func (raw *JSON) UnmarshalJSON(data []byte) error {
	if raw == nil {
		return errors.New("dynamic.RawMessage: UnmarshalJSON on nil pointer")
	}
	*raw = append((*raw)[0:0], data...)
	return nil
}
func (raw JSON) IsObject() bool {
	if len(raw) == 0 {
		return false
	}
	return raw[0] == '{'

}
func (raw JSON) IsMalformed() bool {
	if len(raw) == 0 {
		return true
	}
	switch raw[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		_, err := strconv.ParseFloat(string(raw), 64)
		return err != nil
	case '[':
		if len(raw) == 1 || raw[len(raw)-1] != ']' {
			return true
		}
		var t []interface{}
		err := json.Unmarshal(raw, &t)
		return err == nil
	case '{':
		if len(raw) == 1 || raw[len(raw)-1] != '}' {
			return true
		}
		var t map[string]interface{}
		err := json.Unmarshal(raw, &t)
		return err != nil
	case '"':
		if len(raw) == 1 || raw[len(raw)-1] != '"' {
			return true
		}
		var t string
		err := json.Unmarshal(raw, &t)
		return err != nil
	case 't':
		return !raw.Equal(trueBytes)
	case 'f':
		return !raw.IsFalse()
	case 'n':
		return !raw.IsNull()
	default:
		return true
	}

}

func (raw JSON) IsArray() bool {
	if len(raw) == 0 {
		return false
	}
	return raw[0] == '['
}

func (raw JSON) IsNull() bool {
	return bytes.Equal(raw, Null)
}

// IsBool only reports true if:
//  raw == []byte("true") || raw == []byte("false")
// It does not attempt to parse string values
func (raw JSON) IsBool() bool {
	if len(raw) < 4 {
		return false
	}
	return raw[0] == 't' || raw[0] == 'f'
}

func (raw JSON) IsTrue() bool {
	return bytes.Equal(raw, trueBytes)
}

func (raw JSON) IsFalse() bool {
	return bytes.Equal(raw, falseBytes)
}

func (raw JSON) Equal(data []byte) bool {
	return bytes.Equal(raw, data)
}

// ContainsEscapeRune reports whether the string value of raw contains "\"
// It returns false if raw is not a quoted string.
func (raw JSON) ContainsEscapeRune() bool {
	for i := 0; i < len(raw); i++ {
		if raw[i] == '\\' {
			return true
		}
	}
	return false
}

// UnquotedString trims double quotes from the bytes. It does not parse for
// escaped characters
func (raw JSON) UnquotedString() string {
	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		return string(raw[1 : len(raw)-1])
	}
	return string(raw)
}

// String returns the string representation of the data.
func (raw JSON) String() string {
	if len(raw) == 0 {
		return ""
	}
	return string(raw)
}

func (raw JSON) IsNumber() bool {
	if len(raw) == 0 {
		return false
	}
	switch raw[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return true
	default:
		return false
	}
}
func (raw JSON) IsString() bool {
	if len(raw) == 0 {
		return false
	}
	return raw[0] == '"'
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
