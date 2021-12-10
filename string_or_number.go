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
// checking, use:
//  v := &dynamic.StringOrNumber{}
//  err := v.Set(myType)
func NewStringOrNumber(value interface{}) (StringOrNumber, error) {
	sn := StringOrNumber{
		str:    String{},
		number: Number{},
	}
	err := sn.Set(value)
	return sn, err
}

func NewStringOrNumberPtr(value interface{}) (*StringOrNumber, error) {
	sn, err := NewStringOrNumber(value)
	return &sn, err
}

func (sn StringOrNumber) Value() interface{} {
	switch {
	case sn.number.HasValue():
		return sn.number.Value()
	default:
		return sn.str.String()
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
func (sn *StringOrNumber) Set(value interface{}) error {
	sn.Clear()
	if isNumber(value) {
		return sn.number.Set(value)
	}
	return sn.str.Set(value)

}

func (sn StringOrNumber) MarshalJSON() ([]byte, error) {
	if sn.IsNil() {
		return Null, nil
	}
	if sn.number.HasValue() {
		return sn.number.MarshalJSON()
	}
	return sn.str.MarshalJSON()
}

func (sn *StringOrNumber) UnmarshalJSON(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)
	dec.UseNumber()
	var v interface{}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	return sn.Set(v)
}

// Clear sets the value of StringOrNumber to nil
func (sn *StringOrNumber) Clear() {
	sn.number.Clear()
	sn.str.Clear()
}

func (sn StringOrNumber) String() string {
	switch {
	case sn.number.HasValue():
		return sn.number.String()
	default:
		return sn.str.String()
	}
}

func (sn StringOrNumber) IsNil() bool {
	return sn.str.IsNil() && sn.number.IsNil()
}

// IsString reports whether the value is a string
//
// Note: if you've used any of the other type checks before this that reported
// true, the string may have been cast and is no longer a string.
func (sn StringOrNumber) IsString() bool {
	return sn.str.HasValue()
}

// IsBool reports whether the value is a boolean or a string representation of a
// boolean.
//
// The underlying type of StringOrNumber is altered from a string to a
// bool if the value is successfully parsed.
func (sn *StringOrNumber) IsBool() bool {
	_, ok := sn.Bool()
	return ok
}

// IsEmptyString returns true if sn is nil or an empty string
func (sn *StringOrNumber) IsEmptyString() bool {
	if sn == nil || sn.IsNil() {
		return true
	}
	return sn.String() == ""
}

func (sn *StringOrNumber) Duration() (time.Duration, bool) {
	if sn == nil || !sn.IsString() {
		return time.Duration(0), false
	}
	d, err := time.ParseDuration(sn.str.String())
	if err != nil {
		return time.Duration(0), false
	}
	return d, true
}

func (sn *StringOrNumber) Bool() (bool, bool) {
	if sn.str.HasValue() && !sn.str.IsEmpty() {
		v, err := parseBool(sn.str.String())
		if err != nil {
			return false, false
		}
		sn.str.Clear()
		return *v, true
	}
	return false, false
}

// Time returns the Time value and true. If the original value is a string, Time
// attempts to parse it with the DefaultTimeLayouts or the provided layouts
func (sn *StringOrNumber) Time(layout ...string) (time.Time, bool) {
	if sn.str.HasValue() && !sn.str.IsEmpty() {
		t, err := parseTime(sn.str.String(), layout...)
		if err != nil {
			return time.Time{}, false
		}
		sn.str.value = nil
		return t, true
	}
	return time.Time{}, false
}

func (sn *StringOrNumber) IsTime(layout ...string) bool {
	if _, ok := sn.Time(layout...); ok {
		return ok
	}
	return false
}

// Number returns the underlying value of Number. It could be: int64, uint64,
// float64, or nil. If you need specific types, use their corresponding methods.
func (sn *StringOrNumber) Number() interface{} {
	if sn.number.HasValue() {
		return sn.number.Value()
	}
	if sn.str.HasValue() && !sn.str.IsEmpty() {
		err := sn.number.Parse(sn.str.String())
		if err != nil {
			return nil
		}
		sn.str.Clear()
		return sn.number.Value()
	}
	return nil
}

func (sn *StringOrNumber) Float64() (float64, bool) {
	if sn.number.HasValue() {
		return sn.number.Float64()
	}
	if sn.IsNumber() {
		return sn.number.Float64()
	}
	return 0, false
}

func (sn *StringOrNumber) Float32() (float32, bool) {
	if sn.number.HasValue() {
		return sn.number.Float32()
	}
	if sn.IsNumber() {
		return sn.number.Float32()
	}
	return 0, false
}

func (sn *StringOrNumber) Int64() (int64, bool) {
	if sn.number.HasValue() {
		return sn.number.Int64()
	}
	if sn.IsNumber() {
		return sn.number.Int64()
	}
	return 0, false
}
func (sn *StringOrNumber) Uint64() (uint64, bool) {
	if sn.number.HasValue() {
		return sn.number.Uint64()
	}
	if sn.IsNumber() {
		return sn.number.Uint64()
	}
	return 0, false
}
func (sn *StringOrNumber) Uint() (uint, bool) {
	if sn.number.HasValue() {
		return sn.number.Uint()
	}
	if sn.IsNumber() {
		return sn.number.Uint()
	}
	return 0, false
}
func (sn *StringOrNumber) Uint32() (uint32, bool) {
	if sn.number.HasValue() {
		return sn.number.Uint32()
	}
	if sn.IsNumber() {
		return sn.number.Uint32()
	}
	return 0, false
}
func (sn *StringOrNumber) Uint16() (uint16, bool) {
	if sn.number.HasValue() {
		return sn.number.Uint16()
	}
	if sn.IsNumber() {
		return sn.number.Uint16()
	}
	return 0, false
}
func (sn *StringOrNumber) Uint8() (uint8, bool) {
	if sn.number.HasValue() {
		return sn.number.Uint8()
	}
	if sn.IsNumber() {
		return sn.number.Uint8()
	}
	return 0, false
}
func (sn *StringOrNumber) Int() (int, bool) {
	if sn.number.HasValue() {
		return sn.number.Int()
	}
	if sn.IsNumber() {
		return sn.number.Int()
	}
	return 0, false
}
func (sn *StringOrNumber) Int32() (int32, bool) {
	if sn.number.HasValue() {
		return sn.number.Int32()
	}
	if sn.IsNumber() {
		return sn.number.Int32()
	}
	return 0, false
}
func (sn *StringOrNumber) Int16() (int16, bool) {
	if sn.number.HasValue() {
		return sn.number.Int16()
	}
	if sn.IsNumber() {
		return sn.number.Int16()
	}
	return 0, false
}
func (sn *StringOrNumber) Int8() (int8, bool) {
	if sn.number.HasValue() {
		return sn.number.Int8()
	}
	if sn.IsNumber() {
		return sn.number.Int8()
	}
	return 0, false
}

func (sn *StringOrNumber) IsNumber() bool {
	return sn.Number() != nil
}
