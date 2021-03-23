package dynamic

import (
	"bytes"
	"encoding/json"
	"time"
)

type StringNumberBoolOrTime struct {
	time    Time
	str     String
	number  Number
	boolean Bool
}

func NewStringNumberBoolOrTime(v interface{}) (*StringNumberBoolOrTime, error) {
	snbt := &StringNumberBoolOrTime{
		time:    Time{},
		str:     String{},
		number:  Number{},
		boolean: Bool{},
	}
	err := snbt.Set(v)
	return snbt, err
}

func (snbt StringNumberBoolOrTime) Value() interface{} {
	switch {
	case snbt.time.HasValue():
		return snbt.time.Value()
	case snbt.number.HasValue():
		return snbt.number.Value()
	case snbt.boolean.HasValue():
		return snbt.boolean.value
	default:
		return snbt.str.String()
	}
}

// Set sets the value of StringNumberBoolOrTime to value. The type of value can
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
func (snbt *StringNumberBoolOrTime) Set(value interface{}) error {
	snbt.Clear()
	if isNumber(value) {
		return snbt.number.Set(value)
	}
	if isTime(value) {
		return snbt.time.Set(value)
	}
	if isBool(value) {
		return snbt.boolean.Set(value)
	}
	return snbt.str.Set(value)

}

func (snbt *StringNumberBoolOrTime) MarshalJSON() ([]byte, error) {
	if snbt.IsNil() {
		return json.Marshal(nil)
	}
	if snbt.number.HasValue() {
		return snbt.number.MarshalJSON()
	}
	if snbt.time.HasValue() {
		return snbt.time.MarshalJSON()
	}
	if snbt.boolean.HasValue() {
		return snbt.boolean.MarshalJSON()
	}
	return snbt.str.MarshalJSON()
}

func (snbt *StringNumberBoolOrTime) UnmarshalJSON(data []byte) error {
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

// Clear sets the value of StringNumberBoolOrTime to nil
func (snbt *StringNumberBoolOrTime) Clear() {
	snbt.number.Clear()
	snbt.str.Clear()
	snbt.time.Clear()
	snbt.boolean.Clear()
}

// Format is used only with time. If the value of StringNumberBoolOrTime is anything
// else, the string representation will be returned
func (snbt StringNumberBoolOrTime) Format(layout string) string {
	if snbt.time.HasValue() {
		return snbt.time.Format(layout)
	}
	return snbt.String()
}

func (snbt StringNumberBoolOrTime) String() string {
	switch {
	case snbt.time.HasValue():
		return snbt.time.String()
	case snbt.number.HasValue():
		return snbt.number.String()
	default:
		return snbt.str.String()
	}
}

func (snbt StringNumberBoolOrTime) IsNil() bool {
	return snbt.time.IsNil() && snbt.str.IsNil() && snbt.number.IsNil() && snbt.boolean.IsNil()
}

// IsString reports whether the value is a string
//
// Note: if you've used any of the other type checks before this that reported
// true, the string may have been cast and is no longer a string.
func (snbt StringNumberBoolOrTime) IsString() bool {
	return snbt.str.HasValue()
}

// IsBool reports whether the value is a boolean or a string representation of a
// boolean.
//
// The underlying type of StringNumberBoolOrTime is altered from a string to a
// bool if the value is successfully parsed.
func (snbt *StringNumberBoolOrTime) IsBool() bool {
	_, ok := snbt.Bool()
	return ok
}

func (snbt *StringNumberBoolOrTime) Bool() (bool, bool) {
	if snbt.boolean.HasValue() {
		return snbt.boolean.Bool()
	}
	if snbt.str.HasValue() && !snbt.str.IsEmpty() {
		v, err := parseBool(snbt.str.String())
		if err != nil {
			return false, false
		}
		snbt.str.Clear()
		snbt.boolean.value = v
		return *v, true
	}
	return false, false
}

// Time returns the Time value and true. If the original value is a string, Time
// attempts to parse it with the DefaultTimeLayouts or the provided layouts
func (snbt *StringNumberBoolOrTime) Time(layout ...string) (time.Time, bool) {
	if snbt.time.HasValue() {
		return snbt.time.Time()
	}
	if snbt.str.HasValue() && !snbt.str.IsEmpty() {
		t, err := parseTime(snbt.str.String(), layout...)
		if err != nil {
			return time.Time{}, false
		}
		snbt.str.value = nil
		snbt.time.value = &t
		return t, true
	}
	return time.Time{}, false
}

func (snbt *StringNumberBoolOrTime) IsTime(layout ...string) bool {
	if snbt.time.HasValue() {
		return true
	}
	if _, ok := snbt.Time(layout...); ok {
		return ok
	}
	return false
}

func (snbt *StringNumberBoolOrTime) Number() (interface{}, bool) {
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

func (snbt *StringNumberBoolOrTime) Float() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Float()
	}
	if snbt.IsNumber() {
		return snbt.number.Float()
	}
	return nil, false
}

func (snbt *StringNumberBoolOrTime) Int() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Int()
	}
	if snbt.IsNumber() {
		return snbt.number.Int()
	}
	return nil, false
}
func (snbt *StringNumberBoolOrTime) Uint() (interface{}, bool) {
	if snbt.number.HasValue() {
		return snbt.number.Uint()
	}
	if snbt.IsNumber() {
		return snbt.number.Uint()
	}
	return nil, false
}
func (snbt *StringNumberBoolOrTime) IsNumber() bool {
	_, ok := snbt.Number()
	return ok
}
