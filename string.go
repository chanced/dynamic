package dynamic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type String struct {
	value *string
}

func _() {

}

// NewString returns a new String. Only the first parameter passed in is
// considered.
//
//
// Types
//
// You can set String to any of the following:
//  string, []byte, dynamic.String, fmt.Stringer, []string, *string,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
//
// Warning
//
// This function fails silently. If you pass in an invalid type, the underlying
// value becomes nil. If you want error checking, use:
//  str := dynamic.String{}
//  err := str.Set("value")
// Alternatively, you can check if it is nil:
//  _ = dynamic.NewString().IsNil() // true
func NewString(value ...interface{}) String {
	str := String{}
	if len(value) > 0 {
		_ = str.Set(value[0])
	}
	return str
}

// NewStringPtr returns a pointer to a new String
//
// See NewString for information on valid values and usage
func NewStringPtr(value ...interface{}) *String {
	s := NewString(value...)
	return &s
}

func (s *String) IsNil() bool {
	if s == nil {
		return true
	}
	return s.value == nil
}
func (s String) HasValue() bool {
	return !s.IsEmpty()
}

func (s String) Value() interface{} {
	str := s.value
	if str == nil {
		return nil
	}
	return *str
}

func (s *String) IsEmpty() bool {
	if s.IsNil() {
		return true
	}
	return *s.value == ""
}

func (s String) String() string {
	if s.value == nil {
		return ""
	}
	return *s.value
}
func (s *String) Clear() {
	s.value = nil
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
func (s *String) UnmarshalJSON(data []byte) error {
	s.value = nil
	g := gjson.ParseBytes(data)
	switch g.Type {
	case gjson.Null:
	case gjson.String:
		str := g.String()
		s.value = &str
	default:
		// TODO: really need to do better with errors
		return ErrInvalidValue
	}
	return nil
}

func (s *String) Set(value interface{}) error {
	s.value = nil
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case String:
		return s.Set(v.String())
	case *String:
		return s.Set(*v)
	case string:
		s.value = &v
	case []byte:
		return s.Set(string(v))
	case fmt.Stringer:
		return s.Set(v.String())
	case bool:
		s.Set(strconv.FormatBool(v))
	case int64:
		return s.Set(strconv.FormatInt(v, 10))
	case int32:
		return s.Set(int64(v))
	case int16:
		return s.Set(int64(v))
	case int8:
		return s.Set(int64(v))
	case int:
		return s.Set(int64(v))
	case uint64:
		return s.Set(strconv.FormatUint(v, 10))
	case uint32:
		return s.Set(uint64(v))
	case uint16:
		return s.Set(uint64(v))
	case uint8:
		return s.Set(uint64(v))
	case uint:
		return s.Set(uint64(v))
	case float64:
		return s.Set(strconv.FormatFloat(v, 'f', 0, 64))
	case float32:
		return s.Set(float64(v))
	case complex128:
		return s.Set(strconv.FormatComplex(v, 'f', 0, 128))
	case complex64:
		return s.Set(complex128(v))
	case []string:
		return s.Set(strings.Join(v, ","))
	case *string:
		return s.Set(*v)
	case *bool:
		return s.Set(*v)
	case *int64:
		return s.Set(*v)
	case *int32:
		return s.Set(*v)
	case *int16:
		return s.Set(*v)
	case *int8:
		return s.Set(*v)
	case *int:
		return s.Set(*v)
	case *uint64:
		return s.Set(*v)
	case *uint32:
		return s.Set(*v)
	case *uint16:
		return s.Set(*v)
	case *uint8:
		return s.Set(*v)
	case *uint:
		return s.Set(*v)
	case *float64:
		return s.Set(*v)
	case *float32:
		return s.Set(*v)
	case *complex128:
		return s.Set(*v)
	case *complex64:
		return s.Set(*v)
	default:
		return fmt.Errorf("%w <%T>", ErrInvalidType, value)
	}
	return nil
}
