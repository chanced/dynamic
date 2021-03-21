package dynamic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

var DefaultTimeLayouts = []string{time.RFC3339}
var DefaultTimeLayout = func() string {
	return DefaultTimeLayouts[0]
}

type Time struct {
	value *time.Time
}

// Value returns the underlying value. If it is nil, nil is returned.
func (t Time) Value() interface{} {
	tv := t.value
	if tv == nil {
		return nil
	}
	return *tv
}

// Time returns the underlying value and true if not nil. If it is nil, a zero
// value time is returned and false, indicating as such.
func (t Time) Time() (time.Time, bool) {
	if t.value == nil {
		return time.Time{}, false
	}
	return *t.value, true
}

func (t Time) HasValue() bool {
	return !t.IsNil()
}
func (t Time) IsNil() bool {
	return t.value == nil
}

func (t *Time) Clear() {
	t.Set(nil)
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(t.String())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	t.value = nil
	g := gjson.ParseBytes(data)
	switch g.Type {
	case gjson.Null:
	case gjson.String:
		return t.Parse(g.String())
	default:
		return ErrInvalidValue
	}
	return nil
}

func (t *Time) Set(value interface{}, layout ...string) error {
	t.value = nil
	switch v := value.(type) {
	case time.Time:
		t.value = &v
	case *time.Time:
		tv := *v
		t.value = &tv
	case string:
		return t.Parse(v, layout...)
	case *string:
		return t.Parse(*v, layout...)
	case fmt.Stringer:
		return t.Parse(v.String(), layout...)
	case Time:
		if v.value == nil {
			return nil
		}
		tv := *v.value
		t.value = &tv
	case *Time:
		if v.value == nil {
			return nil
		}
		tv := *v.value
		t.value = &tv
	case nil:
		return nil
	}
	return ErrInvalidValue
}

func isTime(v interface{}) bool {
	switch v.(type) {
	case time.Time, *time.Time, Time, *Time:
		return true
	default:
		return false
	}
}

func (t *Time) Format(layout string) string {
	return t.value.Format(layout)
}

func parseTime(s string, layouts ...string) (time.Time, error) {
	var lastErr error
	if len(layouts) == 0 {
		layouts = DefaultTimeLayouts
	}
	for _, l := range layouts {
		parsed, err := time.Parse(l, s)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func (t *Time) Parse(s string, layouts ...string) error {
	t.value = nil
	tv, err := parseTime(s, layouts...)
	if err != nil {
		return err
	}
	t.value = &tv
	return nil
}

func (t Time) String() string {
	if t.value == nil {
		return ""
	}
	return t.value.Format(DefaultTimeLayout())
}
