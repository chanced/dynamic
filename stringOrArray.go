package dynamic

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrIndexOutOfBounds = errors.New("index out of bounds")
	Done                = errors.New("done")
)

type StringOrArrayOfStrings []string

var typeStringOrArrayOfStrings = reflect.TypeOf(StringOrArrayOfStrings{})

// MarshalJSON satisfies json.Marshaler interface
func (sas StringOrArrayOfStrings) MarshalJSON() ([]byte, error) {
	if sas == nil || len(sas) == 0 {
		return []byte{}, nil
	}
	if len(sas) == 1 {
		return json.Marshal(sas[0])
	}
	return json.Marshal(sas)

}
func (sas *StringOrArrayOfStrings) IsNil() bool {
	return sas == nil || *sas == nil
}

func (sas *StringOrArrayOfStrings) IsEmpty() bool {
	return sas == nil || len(*sas) == 0
}

// Iterate calls fn for each value of sas until an error is returned. If the
// error returned is dynamic.Done, nil is returned to the caller. Otherwise the
// error returned from fn is passed along.
func (sas *StringOrArrayOfStrings) Iterate(fn func(string) error) error {
	if sas == nil {
		return nil
	}
	if sas.IsEmpty() {
		return nil
	}
	for _, v := range *sas {
		err := fn(v)
		if err != nil && err != Done {
			return err
		}
		if err == Done {
			return nil
		}
	}
	return nil
}

// UnmarshalJSON satisfies json.Unmarshaler interface
func (sas *StringOrArrayOfStrings) UnmarshalJSON(data []byte) error {
	b := RawJSON(data)

	if b.IsArray() {
		var v []string
		err := json.Unmarshal(b, &v)
		if err != nil {
			return err
		}
		*sas = v
		return nil
	}
	if b.IsString() {
		var v string
		err := json.Unmarshal(b, &v)
		if err != nil {
			return err
		}
		*sas = append((*sas)[0:0], v)
		return nil
	}
	if b.IsNull() {
		*sas = nil
		return nil
	}
	return &json.UnmarshalTypeError{}

}

func (sas *StringOrArrayOfStrings) Set(v interface{}) error {
	switch t := v.(type) {
	case string:
		*sas = []string{t}
	case []string:
		*sas = t
	case StringOrArrayOfStrings:
		*sas = t
	case *StringOrArrayOfStrings:
		*sas = *t
	case fmt.Stringer:
		*sas = []string{t.String()}
	}
	return ErrInvalidValue
}

func (sas StringOrArrayOfStrings) GetIndex(i int) (string, error) {
	if len(sas) > i || i < 0 {
		return "", ErrIndexOutOfBounds
	}
	return sas[i], nil
}

func (sas *StringOrArrayOfStrings) Add(v string) {
	if sas == nil {
		*sas = StringOrArrayOfStrings{}
	}
	*sas = append(*sas, v)
}

func (sas StringOrArrayOfStrings) IndexOf(v string) int {
	for i, e := range sas {
		if e == v {
			return i
		}
	}
	return -1
}

func (sas *StringOrArrayOfStrings) RemoveIndex(i int) error {
	if i > len(*sas) || i < 0 {
		return ErrIndexOutOfBounds
	}
	if *sas == nil {
		*sas = StringOrArrayOfStrings{}
	}
	*sas = append((*sas)[:i], (*sas)[i+1:]...)
	return nil
}
