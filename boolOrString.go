package dynamic

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// errors
var (
	ErrInvalidValue = errors.New("error: invalid value")
	ErrInvalidType  = errors.New("error: invalid type")
)
var typeBoolOrString = reflect.TypeOf(BoolOrString{})

// BoolOrString is a dynamic type that is a string, bool, or nil.
// It stores the value as a string and decodes json into either a bool or a string
type BoolOrString struct {
	boolean      Bool
	str          String
	encodeToNull bool
}

func (bs *BoolOrString) EncodeToNull() *BoolOrString {
	bs.encodeToNull = true
	return bs
}
func (bs *BoolOrString) EncodeToEmptyString() *BoolOrString {
	bs.encodeToNull = false
	return bs
}
func (bs *BoolOrString) Reference() *BoolOrString {
	return bs
}
func (bs *BoolOrString) Dereference() BoolOrString {
	if bs == nil {
		return BoolOrString{}
	}
	return *bs
}

func (bs BoolOrString) MarshalJSON() ([]byte, error) {
	if bs.IsNil() {
		return Null, nil
	}

	fmt.Println("bs.IsBool():", bs.IsBool(), "for", bs.String())
	if bs.IsBool() {
		b, _ := bs.Bool()
		return json.Marshal(b)
	}
	return json.Marshal(bs.String())
}

func (bs *BoolOrString) UnmarshalJSON(data []byte) error {
	bs.boolean = Bool{}
	bs.str = String{}
	r := RawJSON(data)
	if r.IsNull() {
		return nil
	}
	if r.IsString() {
		var v string
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
		bs.str.Set(v)
		return nil
	}
	if r.IsBool() {
		var v bool
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
		bs.boolean.Set(v)
		return nil
	}
	return &json.UnmarshalTypeError{Value: string(data), Type: typeBoolOrString}
}

// Set sets the BoolOrString's value
//
//
// Valid value types:
//
// You can set String to any of the following:
//  string, []byte, dynamic.String, dynamic.Bool, fmt.Stringer, []string, *string,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
func (bs *BoolOrString) Set(value interface{}) error {
	bs.boolean = Bool{}
	bs.str = String{}
	err := bs.boolean.Set(value)
	if err == nil {
		return nil
	}
	return bs.str.Set(value)
}

func (bs BoolOrString) String() string {
	if !bs.str.IsNil() {
		return bs.str.String()
	}
	if !bs.boolean.IsNil() {
		return bs.boolean.String()
	}
	return ""
}

func (bs *BoolOrString) IsNil() bool {
	if bs == nil {
		return true
	}
	return bs.str.IsNil() && bs.boolean.IsNil()
}

func (bs *BoolOrString) IsEmpty() bool {
	if bs == nil {
		return true
	}
	return bs.boolean.IsNil() && bs.str.IsEmpty()
}

func (bs *BoolOrString) Bool() (value bool, isBool bool) {
	if bs == nil {
		return false, false
	}
	if !bs.IsNil() && !bs.boolean.IsNil() {
		fmt.Println(bs.String(), "!bs.IsNil() && !bs.boolean.IsNil()", *bs.boolean.value)

		return *bs.boolean.value, true
	}
	if !bs.str.IsEmpty() {
		v, err := bs.str.Bool()
		if err == nil {
			return v, true
		}
	}
	return false, false
}

func (bs *BoolOrString) IsBool() bool {
	if bs == nil {
		return false
	}
	_, is := bs.Bool()
	return is
}

// NewBoolOrString returns a new BoolOrString
// It panics if value can not be interpreted as a bool or a string
//
// Valid value types
//
// You can set String to any of the following:
//
//  string, []byte, dynamic.String, dynamic.Bool, fmt.Stringer, []string, *string,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
func NewBoolOrString(value ...interface{}) BoolOrString {
	bs := BoolOrString{}
	if len(value) > 0 {
		err := bs.Set(value[0])
		if err != nil {
			panic(err)
		}
	}
	return bs
}

// NewBoolOrStringPtr returns a pointer to a new BoolOrString
//
// It panics if value can not be interpreted as a bool or a string
// Valid value types:
//
// You can set String to any of the following:
//  string, []byte, dynamic.String, dynamic.Bool, fmt.Stringer, []string, *string,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, complex128, complex64, *float64, *float32, *complex128, *complex64
//  bool, *bool
//  nil
func NewBoolOrStringPtr(value ...interface{}) *BoolOrString {
	bs := NewBoolOrString()
	return bs.Reference()
}
