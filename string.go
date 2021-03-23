package dynamic

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type String struct {
	value *string
}

func (s String) IsNil() bool {
	return s.value == nil
}
func (s String) HasValue() bool {
	return !s.IsNil()
}

func (s String) Value() interface{} {
	str := s.value
	if str == nil {
		return nil
	}
	return *str
}

func (s *String) IsEmpty() bool {

	if s.value == nil {
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
	s.Set(nil)
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
	switch v := value.(type) {
	case String:
		if v.value == nil {
			return nil
		}
		t := v.String()
		s.value = &t
	case *String:
		return s.Set(*v)
	case string:
		s.value = &v
	case *string:
		str := *v
		s.value = &str
	case fmt.Stringer:
		str := v.String()
		s.value = &str
	case nil:
		s.value = nil
	default:
		return fmt.Errorf("%w: expected string or fmt.Stringer, received <%t>", ErrInvalidValue, value)
	}
	return nil
}
