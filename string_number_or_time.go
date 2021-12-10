package dynamic

import (
	"bytes"
	"encoding/json"
	"time"
)

// NewStringNumberOrTime returns a new StringNumberOrTime set to the first
// value, if any.
//
// You can set String to any of the following:
//  string, []byte, dynamic.String, fmt.Stringer, []string, *string,
//  time.Time, *time.Time,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
//
// If you need type checking, use Set
func NewStringNumberOrTime(value interface{}) (StringNumberOrTime, error) {
	snt := StringNumberOrTime{
		time:   Time{},
		str:    String{},
		number: Number{},
	}
	err := snt.Set(value)
	return snt, err
}

// NewStringNumberOrTimePtr returns a pointer to a new NewStringNumberOrTimePtr
//
// See NewNewStringNumberOrTime for info & warnings
func NewStringNumberOrTimePtr(value interface{}) (*StringNumberOrTime, error) {
	snt, err := NewStringNumberOrTime(value)
	return &snt, err
}

// TODO: add format to the instances

type StringNumberOrTime struct {
	time   Time
	str    String
	number Number
}

func (snt StringNumberOrTime) Value() interface{} {
	switch {
	case snt.time.HasValue():
		return snt.time.Value()
	case snt.number.HasValue():
		return snt.number.Value()
	default:
		return snt.str.String()
	}
}

// Set sets the value of StringNumberOrTime to value. The type of value can
// be any of the following:
//
//  string, []byte, dynamic.String, fmt.Stringer, []string, *string,
//  time.Time, *time.Time,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
//
//
// All pointer values are dereferenced.
//
// Set returns an error if value is not one of the aforementioned types
func (snt *StringNumberOrTime) Set(value interface{}) error {
	snt.Clear()
	if value == nil {
		return nil
	}
	if v, ok := value.(*StringNumberOrTime); ok {
		return snt.Set(v.Value())
	}
	if v, ok := value.(StringNumberOrTime); ok {
		return snt.Set(v.Value())
	}
	if isNumber(value) {
		return snt.number.Set(value)
	}
	if isTime(value) {
		return snt.time.Set(value)
	}

	return snt.str.Set(value)

}

func (snt StringNumberOrTime) MarshalJSON() ([]byte, error) {
	if snt.IsNil() {
		return json.Marshal(nil)
	}
	if snt.number.HasValue() {
		return snt.number.MarshalJSON()
	}
	if snt.time.HasValue() {
		return snt.time.MarshalJSON()
	}
	return snt.str.MarshalJSON()
}

func (snt *StringNumberOrTime) UnmarshalJSON(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)
	dec.UseNumber()
	var v interface{}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	return snt.Set(v)
}

// Clear sets the value of StringNumberOrTime to nil
func (snt *StringNumberOrTime) Clear() {
	snt.number.Clear()
	snt.str.Clear()
	snt.time.Clear()
}

// Format is used only with time. If the value of StringNumberOrTime is anything
// else, the string representation will be returned
func (snt StringNumberOrTime) Format(layout string) string {
	if snt.time.HasValue() {
		return snt.time.Format(layout)
	}
	return snt.String()
}

func (snt StringNumberOrTime) String() string {
	switch {
	case snt.time.HasValue():
		return snt.time.String()
	case snt.number.HasValue():
		return snt.number.String()
	default:
		return snt.str.String()
	}
}

func (snt StringNumberOrTime) IsNil() bool {
	return snt.time.IsNil() && snt.str.IsNil() && snt.number.IsNil()
}

// IsString reports whether the value is a string
func (snt StringNumberOrTime) IsString() bool {
	return snt.str.HasValue()
}

// Time returns the Time value and true. If the original value is a string, Time
// attempts to parse it with the DefaultTimeLayouts or the provided layouts
func (snt *StringNumberOrTime) Time(layout ...string) (time.Time, bool) {
	if snt.time.HasValue() {
		return snt.time.Time()
	}
	if snt.str.HasValue() && !snt.str.IsEmpty() {
		t, err := parseTime(snt.str.String(), layout...)
		if err != nil {
			return time.Time{}, false
		}
		snt.str.Clear()
		snt.time.value = &t
		return t, true
	}
	return time.Time{}, false
}

func (snt *StringNumberOrTime) IsTime(layout ...string) bool {
	if snt.time.HasValue() {
		return true
	}
	if _, ok := snt.Time(layout...); ok {
		return ok
	}
	return false
}

// Number returns the underlying value of Number. It could be: int64, uint64,
// float64, or nil. If you need specific types, use their corresponding methods.
func (snt *StringNumberOrTime) Number() interface{} {
	if snt.number.HasValue() {
		return snt.number.Value()
	}
	if snt.str.HasValue() && !snt.str.IsEmpty() {
		err := snt.number.Parse(snt.str.String())
		if err != nil {
			return nil
		}
		snt.str.Clear()
		return snt.number.Value()
	}
	return nil
}

func (snt *StringNumberOrTime) Float64() (float64, bool) {
	if snt.number.HasValue() {
		return snt.number.Float64()
	}
	if snt.IsNumber() {
		return snt.number.Float64()
	}
	return 0, false
}

func (snt *StringNumberOrTime) Float32() (float32, bool) {
	if snt.number.HasValue() {
		return snt.number.Float32()
	}
	if snt.IsNumber() {
		return snt.number.Float32()
	}
	return 0, false
}

func (snt *StringNumberOrTime) Int64() (int64, bool) {
	if snt.number.HasValue() {
		return snt.number.Int64()
	}
	if snt.IsNumber() {
		return snt.number.Int64()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Uint64() (uint64, bool) {
	if snt.number.HasValue() {
		return snt.number.Uint64()
	}
	if snt.IsNumber() {
		return snt.number.Uint64()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Uint() (uint, bool) {
	if snt.number.HasValue() {
		return snt.number.Uint()
	}
	if snt.IsNumber() {
		return snt.number.Uint()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Uint32() (uint32, bool) {
	if snt.number.HasValue() {
		return snt.number.Uint32()
	}
	if snt.IsNumber() {
		return snt.number.Uint32()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Uint16() (uint16, bool) {
	if snt.number.HasValue() {
		return snt.number.Uint16()
	}
	if snt.IsNumber() {
		return snt.number.Uint16()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Uint8() (uint8, bool) {
	if snt.number.HasValue() {
		return snt.number.Uint8()
	}
	if snt.IsNumber() {
		return snt.number.Uint8()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Int() (int, bool) {
	if snt.number.HasValue() {
		return snt.number.Int()
	}
	if snt.IsNumber() {
		return snt.number.Int()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Int32() (int32, bool) {
	if snt.number.HasValue() {
		return snt.number.Int32()
	}
	if snt.IsNumber() {
		return snt.number.Int32()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Int16() (int16, bool) {
	if snt.number.HasValue() {
		return snt.number.Int16()
	}
	if snt.IsNumber() {
		return snt.number.Int16()
	}
	return 0, false
}
func (snt *StringNumberOrTime) Int8() (int8, bool) {
	if snt.number.HasValue() {
		return snt.number.Int8()
	}
	if snt.IsNumber() {
		return snt.number.Int8()
	}
	return 0, false
}

func (snt *StringNumberOrTime) IsNumber() bool {
	return snt.Number() != nil
}

func (snt *StringNumberOrTime) IsNilOrEmpty() bool {
	if snt == nil {
		return true
	}
	if snt.number.IsNil() && snt.time.IsNil() && snt.str.IsNil() {
		return true
	}
	if snt.IsString() {
		return snt.str.IsEmpty()
	}
	return false
}

// IsNilOrZero indiciates whether SNT is nil, empty string, or zero value
func (snt *StringNumberOrTime) IsNilOrZero() bool {
	if snt.IsNilOrEmpty() {
		return true
	}
	if v := snt.Number(); v != nil {
		return v == 0
	}
	if v, ok := snt.Time(); ok {
		return v.IsZero()
	}
	return false
}
