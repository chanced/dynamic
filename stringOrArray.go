package dynamic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

var (
	ErrIndexOutOfBounds = errors.New("index out of bounds")
)

type StringOrArrayOfStrings []string

// MarshalJSON satisfies json.Marshaler interface
func (soas StringOrArrayOfStrings) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if soas == nil || len(soas) == 0 {
		return []byte{}, nil
	}
	var err error
	if len(soas) == 1 {
		err = enc.Encode(soas[1])

	} else {
		err = enc.Encode(soas)
	}

	return buf.Bytes(), err
}

// UnmarshalJSON satisfies json.Unmarshaler interface
func (soas *StringOrArrayOfStrings) UnmarshalJSON(data []byte) error {
	r := gjson.ParseBytes(data)
	ar := r.Array()
	if len(ar) > 0 {
		*soas = StringOrArrayOfStrings{}
	} else {
		return nil
	}
	for _, v := range ar {
		if r.Type == gjson.String {
			*soas = append(*soas, v.String())
			continue
		}
		if r.Type == gjson.Null {
			*soas = append(*soas, "")
			continue
		}
		return fmt.Errorf("can not unmarshal %v into StringOrArrayOfStrings", r.Raw)
	}
	return nil
}

func (soas *StringOrArrayOfStrings) Set(v interface{}) error {
	switch t := v.(type) {
	case string:
		*soas = []string{t}
	case []string:
		*soas = t
	case StringOrArrayOfStrings:
		*soas = t
	case *StringOrArrayOfStrings:
		*soas = *t
	case fmt.Stringer:
		*soas = []string{t.String()}
	}
	return ErrInvalidValue
}

func (soas StringOrArrayOfStrings) GetIndex(i int) (string, error) {
	if len(soas) > i || i < 0 {
		return "", ErrIndexOutOfBounds
	}
	return soas[i], nil
}

func (soas *StringOrArrayOfStrings) Add(v string) {
	if soas == nil {
		*soas = StringOrArrayOfStrings{}
	}
	*soas = append(*soas, v)
}

func (soas StringOrArrayOfStrings) IndexOf(v string) int {
	for i, e := range soas {
		if e == v {
			return i
		}
	}
	return -1
}

func (soas *StringOrArrayOfStrings) RemoveIndex(i int) error {
	if i > len(*soas) || i < 0 {
		return ErrIndexOutOfBounds
	}
	if *soas == nil {
		*soas = StringOrArrayOfStrings{}
	}
	*soas = append((*soas)[:i], (*soas)[i+1:]...)
	return nil
}
