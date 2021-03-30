package dynamic

import (
	"bytes"
	"encoding/json"
	"time"
)

type StringOrNumber struct {
	str    String
	number Number
}

// NewStringOrNumber returns a new StringOrNumber set to the
// first value, if any.
//
//
// String can be set to any of the following:
//  string, []byte, dynamic.String, fmt.Stringer, []string, *string,
//  time.Time, *time.Time,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
//
// Warning
//
// This function panics if value is not an accepted type. If you need error
// checking, use:
//  v := &dynamic.StringOrNumber{}
//  err := v.Set(myType)
func NewStringOrNumber(value interface{}) StringOrNumber {
	snbt := StringOrNumber{
		str:    String{},
		number: Number{},
	}
	err := snbt.Set(value)
	if err != nil {
		panic(err)
	}
	return snbt
}

func NewStringOrNumberPtr(value interface{}) *StringOrNumber {
	snbt := NewStringOrNumber(value)
	return &snbt
}

func (snbt StringOrNumber) Value() interface{} {
	switch {
	case snbt.number.HasValue():
		return snbt.number.Value()
	default:
		return snbt.str.String()
	}
}

// Set sets the value of StringOrNumber to value. The type of value can
// be any of the following:
//  time.Time, *time.Time
//  string, *string
//  json.Number, *json.Number
//  float64, *float64, float32, *float32
//  int, *int, int64, *int64, int32, *int32, int16, *int16, int8, *int8
//  uint, *uint, uint64, *uint64, uint32, *uint32, uint16, *uint16, uint8, *uint8
//  bool, *bool
//  fmt.Stringer
//  nil
//
//
// All pointer values are dereferenced.
//
// Set returns an error if value is not one of the aforementioned types
func (snbt *StringOrNumber) Set(value interface{}) error {
	snbt.Clear()
	if isNumber(value) {
		return snbt.number.Set(value)
	}
	return snbt.str.Set(value)

}

func (snbt StringOrNumber) MarshalJSON() ([]byte, error) {
	if snbt.IsNil() {
		return Null, nil
	}
	if snbt.number.HasValue() {
		return snbt.number.MarshalJSON()
	}
	return snbt.str.MarshalJSON()
}

func (snbt *StringOrNumber) UnmarshalJSON(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)
	dec.UseNumber()
	var v interface{}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	return snbt.Set(v)
}

// Clear sets the value of StringOrNumber to nil
func (snbt *StringOrNumber) Clear() {
	snbt.number.Clear()
	snbt.str.Clear()
}

func (snbt StringOrNumber) String() string {
	switch {
	case snbt.number.HasValue():
		return snbt.number.String()
	default:
		return snbt.str.String()
	}
}

func (snbt StringOrNumber) IsNil() bool {
	return snbt.str.IsNil() && snbt.number.IsNil()
}

// IsString reports whether the value is a string
//
// Note: if you've used any of the other type checks before this that reported
// true, the string may have been cast and is no longer a string.
func (snbt StringOrNumber) IsString() bool {
	return snbt.str.HasValue()
}

// IsBool reports whether the value is a boolean or a string representation of a
// boolean.
//
// The underlying type of StringOrNumber is altered from a string to a
// bool if the value is successfully parsed.
func (snbt *StringOrNumber) IsBool() bool {
	_, ok := snbt.Bool()
	return ok
}

// IsEmptyString returns true if snbt is nil or an empty string
func (snbt *StringOrNumber) IsEmptyString() bool {
	if snbt == nil || snbt.IsNil() {
		return true
	}
	return snbt.String() == ""
}

func (snbt *StringOrNumber) Duration() (time.Duration, bool) {
	if snbt == nil || !snbt.IsString() {
		return time.Duration(0), false
	}
	d, err := time.ParseDuration(snbt.str.String())
	if err != nil {
		return time.Duration(0), false
	}
	return d, true
}

func (snbt *StringOrNumber) Bool() (bool, bool) {
	if snbt.str.HasValue() && !snbt.str.IsEmpty() {
		v, err := parseBool(snbt.str.String())
		if err != nil {
			return false, false
		}
		snbt.str.Clear()
		return *v, true
	}
	return false, false
}

// Time returns the Time value and true. If the original value is a string, Time
// attempts to parse it with the DefaultTimeLayouts or the provided layouts
func (snbt *StringOrNumber) Time(layout ...string) (time.Time, bool) {
	if snbt.str.HasValue() && !snbt.str.IsEmpty() {
		t, err := parseTime(snbt.str.String(), layout...)
		if err != nil {
			return time.Time{}, false
		}
		snbt.str.value = nil
		return t, true
	}
	return time.Time{}, false
}

func (snbt *StringOrNumber) IsTime(layout ...string) bool {
	if _, ok := snbt.Time(layout...); ok {
		return ok
	}
	return false
}

func (snbt *StringOrNumber) Number() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Value(), true
	}
	if snbt.str.HasValue() && !snbt.str.IsEmpty() {
		err := snbt.number.Parse(snbt.str.String())
		if err != nil {
			return nil, false
		}
		snbt.str.Clear()
		return snbt.number.Value(), true
	}
	return nil, false
}

func (snbt *StringOrNumber) Float() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Float()
	}
	if snbt.IsNumber() {
		return snbt.number.Float()
	}
	return nil, false
}

func (snbt *StringOrNumber) Int() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Int()
	}
	if snbt.IsNumber() {
		return snbt.number.Int()
	}
	return nil, false
}
func (snbt *StringOrNumber) Uint() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Uint()
	}
	if snbt.IsNumber() {
		return snbt.number.Uint()
	}
	return nil, false
}
func (snbt *StringOrNumber) IsNumber() bool {
	_, ok := snbt.Number()
	return ok
}
