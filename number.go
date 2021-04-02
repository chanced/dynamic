package dynamic

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var typeNumber = reflect.TypeOf(Number{})

const (
	smallestJSONInt = -9007199254740991
	maxJSONInt      = 9007199254740991
	// smallestJSONFloat = float64(-9007199254740991)
	// maxJSONFloat      = float64(9007199254740991)
)

// NewNumber returns a new Number set to the first, if any, parameters.
//
// You can set String to any of the following:
//  string, *string, json.Number, fmt.Stringer,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, *float64, *float32
func NewNumber(value interface{}) (Number, error) {
	n := Number{}
	err := n.Set(value)
	return n, err
}

// NewNumberPtr returns a pointer to a new Number.
// See NewNumber for information & warnings.
func NewNumberPtr(value interface{}) (*Number, error) {
	n, err := NewNumber(value)

	return &n, err
}

type Number struct {
	intValue   *int64
	uintValue  *uint64
	floatValue *float64
}

func (n *Number) Set(value interface{}) error {
	n.intValue = nil
	n.floatValue = nil
	n.floatValue = nil
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case Number:
		return n.Set(v.Value())
	case *Number:
		return n.Set(n.Value())
	case float32:
		// this is the safest way I can come up with atm.
		//  fl := float32(34.34)
		//  d := float64(34.34)
		//  d == float64(fl) // false

		f, _ := strconv.ParseFloat(strconv.FormatFloat(float64(v), 'f', -1, 32), 32)
		n.floatValue = &f
	case float64:
		n.floatValue = &v
	case int:
		i := int64(v)
		n.intValue = &i
	case uint:
		u := uint64(v)
		n.uintValue = &u
	case int8:
		i := int64(v)
		n.intValue = &i
	case int16:
		i := int64(v)
		n.intValue = &i
	case int32:
		i := int64(v)
		n.intValue = &i
	case int64:
		n.intValue = &v
	case uint8:
		u := uint64(v)
		n.uintValue = &u
	case uint16:
		u := uint64(v)
		n.uintValue = &u
	case uint32:
		u := uint64(v)
		n.uintValue = &u
	case uint64:
		n.uintValue = &v
	case string:
		nv, err := parseNumberFromString(string(v))
		if err != nil {
			return err
		}
		return n.Set(nv)
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return err
		}
		return n.Set(f)
	case []byte:
		return n.Set(string(v))
	case *json.Number:
		return n.Set(*v)
	case fmt.Stringer:
		nv, err := parseNumberFromString(v.String())
		if err != nil {
			return err
		}
		return n.Set(nv)
	case *int64:
		return n.Set(*v)
	case *int32:
		return n.Set(*v)
	case *int16:
		return n.Set(*v)
	case *int8:
		return n.Set(*v)
	case *int:
		return n.Set(*v)
	case *uint64:
		return n.Set(*v)
	case *uint32:
		return n.Set(*v)
	case *uint16:
		return n.Set(*v)
	case *uint8:
		return n.Set(*v)
	case *uint:
		return n.Set(*v)
	case *float64:
		return n.Set(*v)
	case *float32:
		return n.Set(*v)
	default:
		return fmt.Errorf("%w: %T is not a number", ErrInvalidValue, value)
	}
	return nil
}

func (n *Number) Clear() {
	n.floatValue = nil
	n.intValue = nil
	n.uintValue = nil

}

func (n Number) HasValue() bool {
	return !n.IsNil()
}

func (n Number) IsNil() bool {
	return n.floatValue == nil && n.intValue == nil && n.uintValue == nil
}

func (n Number) Bytes() []byte {
	return []byte(n.String())
}
func (n Number) Float32() (float32, bool) {
	if f, ok := n.Float64(); ok && math.MaxFloat32 <= math.Abs(f) {
		v, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', -1, 64), 32)
		return float32(v), true
	}
	return 0, false
}
func (n Number) Float64() (float64, bool) {
	if n.IsNil() {
		return 0, false
	}

	if n.floatValue != nil {
		return *n.floatValue, true
	}
	if n.intValue != nil {
		if *n.intValue == int64(float64(*n.intValue)) {
			return float64(*n.intValue), true
		}
		return 0, false
	}
	if n.uintValue != nil {
		u := *n.uintValue
		if u == uint64(float64(*n.uintValue)) {
			return float64(*n.uintValue), true
		}
		return 0, false
	}
	return 0, false
}

func (n Number) Int64() (int64, bool) {
	if n.intValue != nil {
		return *n.intValue, true
	}

	if n.uintValue != nil {
		u := *n.uintValue
		if u > math.MaxInt64 {
			return 0, false
		}
		return int64(u), true
	}

	if n.floatValue != nil {
		f := *n.floatValue
		if _, frac := math.Modf(*n.floatValue); frac == 0 {
			i := int64(f)
			n.floatValue = nil
			n.intValue = &i
			return i, true
		}
		return 0, false
	}
	return 0, false
}

func (n Number) Int() (int, bool) {
	if i, ok := n.Int64(); ok && int64(int(i)) == i {
		return int(i), true
	}
	return 0, false
}

func (n Number) Int32() (int32, bool) {
	if i, ok := n.Int64(); ok && i <= math.MaxInt32 && i >= math.MinInt32 {
		return int32(i), true
	}
	return 0, false
}
func (n Number) Int16() (int16, bool) {
	if i, ok := n.Int64(); ok && i <= math.MaxInt16 && i >= math.MinInt16 {
		return int16(i), true
	}
	return 0, false
}
func (n Number) Int8() (int8, bool) {
	if i, ok := n.Int64(); ok && i <= math.MaxInt8 && i >= math.MinInt8 {
		return int8(i), true
	}
	return 0, false
}

func (n Number) Uint() (uint, bool) {
	if u, ok := n.Uint64(); ok && uint64(uint(u)) == u {
		return uint(u), true
	}
	return 0, false
}

func (n Number) Uint32() (uint32, bool) {
	if u, ok := n.Uint64(); ok && u <= math.MaxUint32 {
		return uint32(u), true
	}
	return 0, false
}
func (n Number) Uint16() (uint16, bool) {
	if u, ok := n.Uint64(); ok && u <= math.MaxUint16 {
		return uint16(u), true
	}
	return 0, false
}

func (n Number) Uint8() (uint8, bool) {
	if u, ok := n.Uint64(); ok && u <= math.MaxUint8 {
		return uint8(u), true
	}
	return 0, false
}

func (n Number) Uint64() (uint64, bool) {
	if n.uintValue != nil {
		return *n.uintValue, true
	}
	if n.intValue != nil {
		i := *n.intValue
		if i < 0 {
			return 0, false
		}
		return uint64(i), true
	}
	if n.floatValue != nil {
		f := *n.floatValue
		if f < 0 {
			return 0, false
		}
		_, frac := math.Modf(f)
		if f >= 0 && frac == 0 {
			u := uint64(f)
			n.floatValue = nil
			n.uintValue = &u
			return u, true
		}

		return 0, false
	}
	return 0, false
}

func (n *Number) Value() interface{} {
	if n.floatValue != nil {
		return *n.floatValue
	}
	if n.uintValue != nil {
		return *n.uintValue
	}
	if n.intValue != nil {
		return *n.intValue
	}
	return nil
}

func (n *Number) UnmarshalJSON(data []byte) error {
	n.Clear()
	r := JSON(data)
	var err error
	var v interface{}
	switch {
	case r.IsNull():
		return nil
	case r.IsNumber():
		v, err = parseNumberFromString(string(data))

	case r.IsString():
		var str string
		err := json.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		v, err = parseNumberFromString(str)
		if err != nil {
			return err
		}

	default:
		return &json.UnmarshalTypeError{Value: string(data), Type: typeNumber}
	}
	if err != nil {
		return &json.UnmarshalTypeError{Value: string(data), Type: typeNumber}
	}
	err = n.Set(v)
	if err != nil {
		return &json.UnmarshalTypeError{Value: string(data), Type: typeNumber}
	}

	return nil
}

func (n Number) String() string {
	if n.floatValue != nil {
		return strconv.FormatFloat(*n.floatValue, 'f', -1, 64)
	}
	if n.uintValue != nil {
		return strconv.FormatUint(*n.uintValue, 10)
	}
	if n.intValue != nil {

		return strconv.FormatInt(*n.intValue, 10)
	}
	return ""
}

func (n *Number) Parse(s string) error {
	n.Clear()
	v, err := parseNumberFromString(s)
	if err != nil {
		return err
	}
	switch t := v.(type) {
	case float64:
		n.floatValue = &t
	case int64:
		n.intValue = &t
	case uint64:
		n.uintValue = &t
	}
	return nil
}

func (n Number) MarshalJSON() ([]byte, error) {
	if n.uintValue != nil {
		u := *n.uintValue
		if maxJSONInt >= u {
			return json.Marshal(u)
		} else {
			return json.Marshal(strconv.FormatUint(u, 10))
		}
	}
	if n.intValue != nil {
		i := *n.intValue
		if maxJSONInt >= i && i >= smallestJSONInt {
			return json.Marshal(i)
		} else {
			return []byte(strconv.FormatInt(i, 10)), nil
		}
	}
	if n.floatValue != nil {
		return json.Marshal(*n.floatValue)
	}
	return json.Marshal(nil)
}

func isNumber(v interface{}) bool {
	switch v.(type) {
	case Number, *Number,
		uint, *uint, uint64, *uint64, uint32, *uint32, uint16, *uint16, uint8, *uint8,
		int, *int, int64, *int64, int32, *int32, int16, *int16, int8, *int8,
		float64, *float64, float32, *float32,
		json.Number, *json.Number:
		return true
	default:
		return false
	}
}

func parseUint(s string) (uint64, bool) {
	u, err := strconv.ParseUint(s, 0, 64)
	return u, err == nil
}

func parseInt(s string) (int64, bool) {
	i, err := strconv.ParseInt(s, 0, 64)
	return i, err == nil
}

func parseFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}

func parseNumberFromString(s string) (interface{}, error) {
	if len(s) == 0 {
		return nil, nil
	}
	if u, ok := parseUint(s); ok {
		return u, nil
	}
	if i, ok := parseInt(s); ok {
		return i, nil
	}
	if f, ok := parseFloat(s); ok {
		return f, nil
	}

	return nil, fmt.Errorf("%w: \"%s\" is not a number", ErrInvalidValue, s)
}
