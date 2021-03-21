package dynamic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// errors
var (
	ErrInvalidValue = errors.New("invalid value")
)

// TODO: use Bool here

// BoolOrString is a dynamic type that is either a string or a bool for json encoding.
// It stores the value as a string and decodes json into either a bool or a string
type BoolOrString string

// MarshalJSON satisfies json.Marshaler interface
func (bos BoolOrString) MarshalJSON() ([]byte, error) {

	v := strings.ToLower(string(bos))
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	var err error
	switch v {
	case "true":
		err = enc.Encode(true)
		return buf.Bytes(), err
	case "false":
		err = enc.Encode(false)
	default:
		err = enc.Encode(string(bos))
	}
	return buf.Bytes(), err
}

// UnmarshalJSON satisfies json.Unmarshaler interface
func (bos *BoolOrString) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	t := r.Type
	switch t {
	case gjson.False:
		*bos = "false"
	case gjson.True:
		*bos = "true"
	case gjson.String:
		*bos = BoolOrString(r.String())
	case gjson.Null:
		*bos = ""
	default:
		return fmt.Errorf("can not unmarshal %v into BoolOrString", r.Num)
	}
	return nil
}

// Set sets the value of BoolOrString to v. Allowed types are: bool, string, or fmt.Stringer
func (bos *BoolOrString) Set(v interface{}) error {
	switch t := v.(type) {
	case string:
		str := strings.ToLower(t)
		if str == "true" || str == "false" {
			*bos = BoolOrString(str)
			return nil
		}
		*bos = BoolOrString(t)
		return nil
	case bool:
		if t {
			*bos = "true"
		} else {
			*bos = "false"
		}
		return nil
	case BoolOrString:
		*bos = t
	case *BoolOrString:
		*bos = *t
	case fmt.Stringer:
		return bos.Set(t.String())
	}
	return fmt.Errorf("%w for BoolOrString: expected bool, string, or fmt.Stringer", ErrInvalidValue)
}

// IsBool returns true if it is a bool or false if it is a string
func (bos BoolOrString) IsBool() bool {
	return bos.Bool() != nil
}

// IsString returns true if it is a string or false if it is a boolean
func (bos BoolOrString) IsString() bool {
	return !bos.IsBool()
}

// Bool returns a pointer boolean. If the value is a boolean,
// it returns true or false otherwise nil is returned
func (bos BoolOrString) Bool() *bool {
	str := strings.ToLower(string(bos))
	if str == "true" {
		t := true
		return &t
	} else if str == "false" {
		f := false
		return &f
	}
	return nil
}

// IsTrue reports whether or not the value is a boolean and "true"
func (bos BoolOrString) IsTrue() bool {
	return bos.Bool() != nil && *bos.Bool()
}

// IsFalse reports whether or not the value is a boolean and "false"
func (bos BoolOrString) IsFalse() bool {
	return bos.Bool() != nil && !*bos.Bool()
}

func (bos BoolOrString) String() string {
	return string(bos)
}
