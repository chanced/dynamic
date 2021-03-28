package dynamic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

var boolType = reflect.TypeOf(Bool{})

// NewBool returns a new Bool value initialized to the first, if any, value
// passed in.
//
// You can set String to any of the following:
//  bool, dynamic.Bool, *bool, *dynamic.Bool
//  string, []byte, fmt.Stringer, *string
//  nil
//
func NewBool(value interface{}) (Bool, error) {
	b := Bool{}

	err := b.Set(value)

	return b, err
}

// NewBoolPtr returns a pointer to a new Bool
//
// See NewBool for valid options, usage and warnings
func NewBoolPtr(value interface{}) (*Bool, error) {
	b, err := NewBool(value)
	return &b, err
}

type Bool struct {
	value *bool
}

var (
	falseValue = false
	trueValue  = true
)
var (
	False = Bool{value: &falseValue}
	True  = Bool{value: &trueValue}
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

// var ErrInvalidBool = errors.New("invalid Bool")

func (b Bool) Equal(value interface{}) bool {
	if b.value == nil {
		return value == nil
	}
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return *b.value == v
	case string:
	case Bool:
		return *b.value == *v.value
	case *Bool:
	case *bool:
		return *b.value == *v

	}

	return false
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b.value == nil {
		return Null, nil
	}
	return json.Marshal(*b.value)
}
func (b *Bool) UnmarshalJSON(data []byte) error {
	b.value = nil
	r := JSON(data)
	switch {
	case r.IsTrue():
		t := true
		b.value = &t
	case r.IsFalse():
		f := false
		b.value = &f
	case r.IsNull():
		return nil
	case r.IsString():
		if r.String() == "" {
			return nil
		}
		v, err := strconv.ParseBool(r.String())
		if err != nil {
			return &json.UnmarshalTypeError{Value: r.String(), Type: boolType}
		}
		b.value = &v
	default:
		return &json.UnmarshalTypeError{Value: r.String(), Type: boolType}
	}
	return nil
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
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case bool:
		b.value = &v
	case Bool:
		b.value = v.value
	case string:
		return b.Parse(v)
	case []byte:
		return b.Set(string(v))
	case *bool:
		b.value = v
	case *Bool:
		return b.Set(v.Value())
	case *string:
		return b.Parse(*v)
	default:
		return fmt.Errorf("%w type: %T", ErrInvalidType, value)
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
func (b *Bool) IsNull() bool {
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
