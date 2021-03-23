package dynamic

import (
	"bytes"
	"encoding/json"
	"time"
)

func NewStringNumberOrTime(v interface{}) (*StringNumberOrTime, error) {
	snt := &StringNumberOrTime{
		time:   Time{},
		str:    String{},
		number: Number{},
	}
	err := snt.Set(v)
	return snt, err
}

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
//  time.Time, *time.Time
//  string, *string
//  json.Number, *json.Number
//  float64, *float64, float32, *float32
//  int, *int, int64, *int64, int32, *int32, int16, *int16, int8, *int8
//  uint, *uint, uint64, *uint64, uint32, *uint32, uint16, *uint16, uint8, *uint8
//  fmt.Stringer
//  nil
//
//
// All pointer values are dereferenced.
//
// Set returns an error if value is not one of the aforementioned types
func (snt *StringNumberOrTime) Set(value interface{}) error {
	snt.Clear()
	if v, ok := value.(*StringNumberOrTime); ok {
		return snt.Set(v.Value())
	}
	if v, ok := value.(*StringNumberOrTime); ok {
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

func (snt *StringNumberOrTime) MarshalJSON() ([]byte, error) {
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

func (snt *StringNumberOrTime) Number() (interface{}, bool) {
	if snt.number.HasValue() {
		return snt.number.Value(), true
	}
	if snt.str.HasValue() && !snt.str.IsEmpty() {
		err := snt.number.Parse(snt.str.String())
		if err != nil {
			return nil, false
		}
		snt.str.Clear()
		return snt.number.Value(), true
	}
	return nil, false
}

func (snt *StringNumberOrTime) Float() (interface{}, bool) {
	if snt.number.HasValue() {
		return snt.number.Float()
	}
	if snt.IsNumber() {
		return snt.number.Float()
	}
	return nil, false
}

func (snt *StringNumberOrTime) Int() (interface{}, bool) {
	if snt.number.HasValue() {
		return snt.number.Int()
	}
	if snt.IsNumber() {
		return snt.number.Int()
	}
	return nil, false
}
func (snt *StringNumberOrTime) Uint() (interface{}, bool) {
	if snt.number.HasValue() {
		return snt.number.Uint()
	}
	if snt.IsNumber() {
		return snt.number.Uint()
	}
	return nil, false
}
func (snt *StringNumberOrTime) IsNumber() bool {
	_, ok := snt.Number()
	return ok
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
	if v, ok := snt.Number(); ok {
		return v == 0
	}
	if v, ok := snt.Time(); ok {
		return v.IsZero()
	}
	return false
}
