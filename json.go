package dynamic

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

type RawJSON []byte

var (
	trueBytes  = []byte("true")
	falseBytes = []byte("false")
)

func (r RawJSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}
func (r *RawJSON) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errors.New("dynamic.RawMessage: UnmarshalJSON on nil pointer")
	}
	*r = append((*r)[0:0], data...)
	return nil
}
func (r RawJSON) IsObject() bool {
	if len(r) == 0 {
		return false
	}
	return r[0] == '{'

}
func (r RawJSON) IsMalformed() bool {
	if len(r) == 0 {
		return true
	}
	switch r[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		_, err := strconv.ParseFloat(string(r), 64)
		return err != nil
	case '[':
		if len(r) == 1 || r[len(r)-1] != ']' {
			return true
		}
		var t []interface{}
		err := json.Unmarshal(r, &t)
		return err == nil
	case '{':
		if len(r) == 1 || r[len(r)-1] != '}' {
			return true
		}
		var t map[string]interface{}
		err := json.Unmarshal(r, &t)
		return err != nil
	case '"':
		if len(r) == 1 || r[len(r)-1] != '"' {
			return true
		}
		var t string
		err := json.Unmarshal(r, &t)
		return err != nil
	case 't':
		return !r.Equal(trueBytes)
	case 'f':
		return !r.IsFalse()
	case 'n':
		return !r.IsNull()
	default:
		return true
	}

}

func (r RawJSON) IsArray() bool {
	if len(r) == 0 {
		return false
	}
	return r[0] == '['
}

func (r RawJSON) IsNull() bool {
	return bytes.Equal(r, Null)
}

// IsBool only reports true if:
//  r == []byte("true") || r == []byte("false")
// It does not attempt to parse string values
func (r RawJSON) IsBool() bool {
	if len(r) < 4 {
		return false
	}
	return r[0] == 't' || r[0] == 'f'
}

func (r RawJSON) IsTrue() bool {
	return bytes.Equal(r, trueBytes)
}

func (r RawJSON) IsFalse() bool {
	return bytes.Equal(r, falseBytes)
}

func (r RawJSON) Equal(data []byte) bool {
	return bytes.Equal(r, data)
}

// ContainsEscapeRune reports whether the string value of r contains "\"
// It returns false if r is not a quoted string.
func (r RawJSON) ContainsEscapeRune() bool {
	for i := 0; i < len(r); i++ {
		if r[i] == '\\' {
			return true
		}
	}
	return false
}

// UnquotedString trims double quotes from the bytes. It does not parse for
// escaped characters
func (r RawJSON) UnquotedString() string {
	if r[0] == '"' && r[len(r)-1] == '"' {
		return string(r[1 : len(r)-1])
	}
	return string(r)
}

// String returns the string representation of the data.
func (r RawJSON) String() string {
	if len(r) == 0 {
		return ""
	}
	return string(r)
}

func (r RawJSON) IsNumber() bool {
	if len(r) == 0 {
		return false
	}
	switch r[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return true
	default:
		return false
	}
}
func (r RawJSON) IsString() bool {
	if len(r) == 0 {
		return false
	}
	return r[0] == '"'
}
