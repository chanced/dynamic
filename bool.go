package dynamic

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Bool struct {
	value *bool
}

var (
	False = Bool{value: func() *bool {
		v := false
		return &v
	}(),
	}
	True = Bool{value: func() *bool {
		v := true
		return &v
	}(),
	}
)

func (b Bool) Value() interface{} {
	if b.value == nil {
		return nil
	}
	return *b.value
}

func (b Bool) HasValue() bool {
	return !b.IsNil()
}

var ErrInvalidBool = errors.New("invalid Bool")

func (b Bool) Equal(value interface{}) bool {
	if b.value == nil {
		return value == nil
	}
	switch v := value.(type) {
	case bool:
		return *b.value == v
	case *bool:
		return *b.value == *v
	case string:
		s := strings.ToLower(v)
		if *b.value {
			return s == "true"
		}
		return s == "false"
	case *string:
		s := strings.ToLower(*v)
		if *b.value {
			return s == "true"
		}
		return s == "false"
	case Bool:
		return *b.value == *v.value
	case *Bool:
		return *b.value == *v.value
	}

	return false
}
func (b Bool) MarshalJSON() ([]byte, error) {
	if b.value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(*b.value)
}
func (b *Bool) UnmarshalJSON(data []byte) error {
	res := gjson.ParseBytes(data)
	switch res.Type {
	case gjson.False:
		*b = False
		return nil
	case gjson.True:
		*b = True
		return nil
	case gjson.String:
		return b.Parse(string(data))
	case gjson.Null:
		return nil
	}
	return ErrInvalidBool
}

func parseBool(str string) (*bool, error) {
	if str == "" {
		return nil, nil
	}
	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (b *Bool) Bool() (bool, bool) {
	if b.value == nil {
		return false, false
	}
	return *b.value, true
}

func (b *Bool) Parse(str string) error {
	b.value = nil
	v, err := parseBool(str)
	if err != nil {
		return err
	}
	b.value = v
	return nil
}

func (b *Bool) Set(value interface{}) error {
	b.value = nil
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		return b.Set(string(v))
	case *bool:
		b.value = v
	case bool:
		b.value = &v
	case string:
		return b.Parse(v)
	case *string:
		return b.Parse(*v)
	case Bool:
		b.value = v.value
	case *Bool:
		if v == nil {
			b.value = nil
			return nil
		}
		return b.Set(v.Value())
	default:
		return fmt.Errorf("%w type: %t", ErrInvalidBool, value)
	}
	return nil
}
func (b *Bool) SetValue(v bool) {
	if b == nil {
		*b = Bool{}
	}
	b.value = &v
}

func (b *Bool) Clear() {
	b.value = nil
}

func (b *Bool) String() string {
	if b.value == nil {
		return ""
	}
	if b.value == nil {
		return ""
	}
	if *b.value {
		return "true"
	}
	return "false"
}

func (b *Bool) IsNil() bool {
	return b == nil || b.value == nil
}

func (b *Bool) IsTrue() bool {
	if b.IsNil() {
		return false
	}
	return *b.value
}

func (b *Bool) IsFalse() bool {
	if b.IsNil() {
		return false
	}
	return !*b.value
}

func isBool(value interface{}) bool {
	if _, ok := value.(bool); ok {
		return true
	}
	if _, ok := value.(*bool); ok {
		return true
	}
	return false
}
