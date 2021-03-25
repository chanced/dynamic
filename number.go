package dynamic

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/tidwall/gjson"
)

const smallestJSONInt = -9007199254740991
const maxJSONInt = 9007199254740991
const smallestJSONFloat = float64(-9007199254740991)
const maxJSONFloat = float64(9007199254740991)

// NewNumber returns a new Number set to the first, if any, parameters.
//
// Types
//
// You can set String to any of the following:
//  string, *string, json.Number, fmt.Stringer,
//  int, int64, int32, int16, int8, *int, *int64, *int32, *int16, *int8,
//  uint, uint64, uint32, uint16, uint8, *uint, *uint64, *uint32, *uint16, *uint8
//  float64, float32, *float64, *float32
//
//
// Warning
//
// This function fails silently. If you pass in an invalid type, the underlying
// value becomes nil. If you want error checking, use:
func NewNumber(value ...interface{}) Number {
	n := Number{}
	if len(value) > 0 {
		_ = n.Set(value[0])
	}
	return n
}

// NewNumberPtr returns a pointer to a new Number.
// See NewNumber for information & warnings.
func NewNumberPtr(value ...interface{}) *Number {
	n := NewNumber(value...)
	return &n
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
		f := float64(v)
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

func (n Number) Float() (float64, bool) {
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

func (n Number) Int() (int64, bool) {
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

func (n Number) Uint() (uint64, bool) {
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
	g := gjson.ParseBytes(data)
	switch g.Type {
	case gjson.Null:
	case gjson.Number:
		f := g.Float()
		return n.Set(f)
	case gjson.String:
		return n.Set(g.String())
	default:
		// TODO: better errors
		return ErrInvalidValue
	}
	return nil
}

func (n Number) String() string {
	if n.floatValue != nil {
		return strconv.FormatFloat(*n.floatValue, 'f', -1, 64)
	}
	if n.uintValue != nil {
		return strconv.FormatUint(*n.uintValue, 64)
	}
	if n.intValue != nil {
		return strconv.FormatInt(*n.intValue, 64)
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
		if maxJSONFloat >= *n.floatValue && *n.floatValue >= smallestJSONFloat {
			return json.Marshal(*n.floatValue)
		}
		return json.Marshal(strconv.FormatFloat(*n.floatValue, 'f', -1, 64))
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
