package dynamic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type String struct {
	value                  *string
	encodeNilAsEmptyString bool
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
func (s *String) Reference() *String {
	return s
}
func (s *String) Dereference() String {
	if s == nil {
		return String{}
	}
	return *s
}

func (s *String) IsNil() bool {
	if s == nil {
		return true
	}
	return s.value == nil
}

// IsNull reports whether s is nil or equal to "null"
func (s *String) IsNull() bool {
	return s == nil || s.IsNil() || s.Equal("null")
}
func (s *String) Len() int {
	if s == nil || s.IsEmpty() {
		return 0
	}
	return len(*s.value)
}

func (s *String) NewReader() *strings.Reader {
	if s == nil || s.IsEmpty() {
		return strings.NewReader("")
	}
	return strings.NewReader(*s.value)
}
func (s *String) ContainsRune(r rune) bool {
	if s == nil || s.IsEmpty() {
		return false
	}
	return strings.ContainsRune(*s.value, r)
}

// Count counts the number of non-overlapping instances of substr in s.
//
// If substr is an empty string, Count returns 1 + the number of Unicode code
// points in s.
func (s *String) Count(value interface{}) int {
	if s == nil || s.IsEmpty() {
		return 0
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return 0
	}
	return strings.Count(s.String(), *str)
}

// Map returns a copy of *String s with all its characters modified according
// to the mapping function. If mapping returns a negative value, the character
// is dropped from the string with no replacement.
func (s *String) Map(mapping func(rune) rune) *String {
	if s == nil {
		return &String{}
	}
	if s.IsEmpty() {
		return s
	}
	return NewStringPtr(strings.Map(mapping, s.String()))

}

// Replace returns a copy of the string s with the first n non-overlapping
// instances of old replaced by new. If old is empty, it matches at the
// beginning of the string and after each UTF-8 sequence, yielding up to k+1
// replacements for a k-rune string. If n < 0, there is no limit on the number
// of replacements.
//
// Replace panics if old or new can not be parsed as strings (see Set)
func (s *String) Replace(old interface{}, new interface{}, n int) *String {
	if s == nil {
		return &String{}
	}
	if s.IsEmpty() {
		return s
	}

	oldPtr, err := formatString(old)
	if err != nil {
		panic(err)
	}
	var oldStr string
	if oldPtr != nil {
		oldStr = *oldPtr
	}
	var newStr string
	newPtr, err := formatString(new)
	if err != nil {
		panic(err)
	}
	if newPtr != nil {
		newStr = *newPtr
	}
	return NewStringPtr(strings.Replace(s.String(), oldStr, newStr, n))
}

// ReplaceAll returns a copy of the string s with all non-overlapping instances
// of old replaced by new. If old is empty, it matches at the beginning of the
// string and after each UTF-8 sequence, yielding up to k+1 replacements for a
// k-rune string.
//
// Replace panics if old or new can not be parsed as strings (see Set)
func (s *String) ReplaceAll(old interface{}, new interface{}) *String {
	if s == nil {
		return &String{}
	}
	return s.Replace(old, new, -1)

}

// SplitN slices s into substrings separated by sep and returns a slice of the
// substrings between those separators.
//
// The count determines the number of substrings to return:
func (s *String) SplitN(sep interface{}, n int) []string {
	if s == nil || s.IsEmpty() {
		return []string{}
	}

	ptr, err := formatString(sep)
	if err != nil {
		panic(err)
	}
	if ptr == nil || *ptr == "" {
		return []string{s.String()}
	}
	return strings.SplitN(s.String(), *ptr, n)
}

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// If s does not contain sep and sep is not empty, Split returns a slice of
// length 1 whose only element is s.
//
// If sep is empty, Split splits after each UTF-8 sequence. If both s and sep
// are empty, Split returns an empty slice.
//
// It is equivalent to SplitN with a count of -1.
func (s *String) Split(sep interface{}) []string {
	if s == nil || s.IsEmpty() {
		return []string{}
	}
	return s.SplitN(sep, -1)
}

// SplitAfter slices s into all substrings after each instance of sep and
// returns a slice of those substrings.
//
// If s does not contain sep and sep is not empty, SplitAfter returns a slice of
// length 1 whose only element is s.
//
// If sep is empty, SplitAfter splits after each UTF-8 sequence. If both s and
// sep are empty, SplitAfter returns an empty slice.
//
// It is equivalent to SplitAfterN with a count of -1.
func (s *String) SplitAfter(sep interface{}) []string {
	if s == nil || s.IsEmpty() {
		return []string{}
	}
	return s.SplitAfterN(sep, -1)
}
func (s *String) SplitAfterN(sep interface{}, n int) []string {
	if s == nil || s.IsEmpty() {
		return []string{}
	}
	ptr, err := formatString(sep)
	if err != nil {
		panic(err)
	}
	if ptr == nil || *ptr == "" {
		return []string{s.String()}
	}
	return strings.SplitAfterN(s.String(), *ptr, n)
}

// Title returns a copy of the *String s with all Unicode letters that begin words
// mapped to their Unicode title case.
//
// BUG(rsc): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
func (s *String) Title() *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.Title(s.String()))

}

// ToLower returns a new *String with all Unicode letters mapped to their lower case.
func (s *String) ToLower() *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.ToLower(s.String()))
}

// ToLowerSpecial returns a copy of the *String s with all Unicode letters
// mapped to their lower case using the case mapping specified by c.
func (s *String) ToLowerSpecial(c unicode.SpecialCase) *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.ToLowerSpecial(c, s.String()))
}

// ToTitle returns a copy of the *String s with all Unicode letters mapped to
// their Unicode title case.
func (s *String) ToTitle() *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.ToTitle(s.String()))
}

// ToTitleSpecial returns a copy of the *String s with all Unicode letters
// mapped to their Unicode title case, giving priority to the special casing
// rules.
func (s *String) ToTitleSpecial(c unicode.SpecialCase) *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.ToTitleSpecial(c, s.String()))

}

// ToUpper returns s with all Unicode letters mapped to their upper case.
func (s *String) ToUpper() *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.ToUpper(s.String()))
}

// ToUpperSpecial returns a copy of the *String s with all Unicode letters
// mapped to their upper case using the case mapping specified by c.
func (s *String) ToUpperSpecial(c unicode.SpecialCase) *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(strings.ToUpperSpecial(c, s.String()))

}

// ToValidUTF8 returns a copy of the string s with each run of invalid UTF-8
// byte sequences replaced by the replacement string, which may be empty.
//
// It panics if replacement cannot be formatted into a string. See Set.
func (s *String) ToValidUTF8(replacement interface{}) *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	repl, err := formatString(replacement)
	if err != nil {
		panic(err)
	}
	if repl == nil {
		return NewStringPtr(s)
	}
	return NewStringPtr(strings.ToValidUTF8(s.String(), *repl))
}

func (s *String) Copy() *String {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return &String{}
	}
	return NewStringPtr(s.String())
}

// ContainsAny reports whether any Unicode code points in chars are within s.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set
func (s *String) ContainsAny(value interface{}) bool {
	if s == nil || s.IsNil() {
		return false
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil || *str == "" {
		return false
	}
	return strings.ContainsAny(s.String(), *str)
}

//Compare returns an integer comparing two strings lexicographically.
//
// The result will be 0 if s==value, -1 if s < value, and +1 if s > value.
//
// It panics if value can not be interpreted as a string or formatted as such. See Set
func (s *String) Compare(value interface{}) (res int) {
	if s == nil || s.IsNil() {
		return -1
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		empty := ""
		str = &empty
	}
	return strings.Compare(s.String(), *str)
}

// Contains reports whether the formatted value is within s.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set
func (s *String) Contains(value interface{}) bool {
	if s == nil || s.IsNil() {
		return false
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return false
	}
	return strings.Contains(s.String(), *str)

}

// Fields splits the string s around each instance of one or more consecutive
// white space characters, as defined by unicode.IsSpace, returning a slice of
// substrings of s or an empty slice if s contains only white space.
func (s *String) Fields() []string {
	if s == nil || s.IsEmpty() {
		return []string{}
	}
	return strings.Fields(*s.value)
}
func (s *String) FieldsFunc(f func(rune) bool) []string {
	if s == nil || s.IsEmpty() {
		return []string{}
	}

	return strings.FieldsFunc(s.String(), f)
}

// HasPrefix tests whether the string s begins with prefix.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set
func (s *String) HasPrefix(value interface{}, timeLayouts ...string) bool {
	if s == nil || s.IsNil() {
		return false
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return false
	}
	return strings.HasPrefix(s.String(), *str)

}

// HasSuffix tests whether the string s ends with suffix.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set
func (s *String) HasSuffix(value interface{}) bool {
	if s == nil || s.IsNil() {
		return false
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return false
	}
	return strings.HasSuffix(s.String(), *str)

}

// LastIndex returns the index of the last instance of substr in s, or -1 if
// substr is not present in s.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set
func (s *String) LastIndex(value interface{}) (res int) {
	if s == nil || s.IsNil() {
		return -1
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return -1
	}
	return strings.LastIndex(s.String(), *str)
}

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set.
func (s *String) IndexAny(value interface{}) (res int) {
	if s == nil || s.IsNil() {
		return -1
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return -1
	}
	return strings.IndexAny(s.String(), *str)
}

// Index returns the index of the first instance of substr in s, or -1 if substr
// is not present in s.
//
// It panics if value can not be interpreted as a string or formatted as such.
// See Set
func (s *String) Index(value interface{}) (res int) {
	if s == nil || s.IsNil() {
		return -1
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return -1
	}
	return strings.Index(s.String(), *str)
}

// LastIndexFunc returns the index into s of the last Unicode code point
// satisfying f(c), or -1 if none do.
func (s *String) LastIndexFunc(fn func(r rune) bool) int {
	if s == nil || s.IsEmpty() {
		return -1
	}

	return strings.LastIndexFunc(s.String(), fn)
}

// IndexFunc returns the index into s of the first Unicode code point satisfying
// f(c), or -1 if none do.
func (s *String) IndexFunc(fn func(r rune) bool) int {
	if s == nil || s.IsEmpty() {
		return -1
	}
	return strings.IndexFunc(s.String(), fn)
}

// EqualFold reports whether s and t, interpreted as UTF-8 strings, are equal
// under Unicode case-folding, which is a more general form of
// case-insensitivity.
//
// It panics if value can not be interpreted or formatted as a string. See Set
func (s *String) EqualFold(value interface{}) bool {
	if s == nil || s.value == nil {
		return value == nil
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	if str == nil {
		return false
	}
	return strings.EqualFold(s.String(), *str)
}

// Equal reports whether the formatted value is equal to the underlying value
// of *String s.
//
// It panics if value can not be formatted as a string. See Set
func (s *String) Equal(value interface{}) bool {
	if s == nil || s.value == nil {
		return value == nil
	}
	str, err := formatString(value)
	if err != nil {
		panic(err)
	}
	return *s.value == *str
}

// IndexByte returns the index of the first instance of c in s, or -1 if c is
// not present in s.
func (s *String) IndexByte(c byte) int {
	if s == nil || s.IsEmpty() {
		return -1
	}
	return strings.IndexByte(*s.value, c)

}

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is
// not present in s.
func (s *String) LastIndexByte(c byte) int {
	if s == nil || s.IsEmpty() {
		return -1
	}
	return strings.LastIndexByte(*s.value, c)

}

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s. If r is utf8.RuneError, it returns the
// first instance of any invalid UTF-8 byte sequence.
func (s *String) IndexRune(r rune) int {
	if s == nil || s.IsEmpty() {
		return -1
	}
	return strings.IndexRune(*s.value, r)
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

func (s *String) EncodeNilToNull() {
	s.encodeNilAsEmptyString = false
}

func (s *String) EncodeNilToEmptyString() {
	s.encodeNilAsEmptyString = true
}

func (s *String) IsEmpty() bool {
	if s.IsNil() {
		return true
	}
	return *s.value == ""
}

func (s *String) String() string {
	if s == nil || s.value == nil {
		return ""
	}
	return *s.value
}
func (s *String) Clear() {
	s.value = nil
}

func (s String) MarshalJSON() ([]byte, error) {
	if s.IsNil() && !s.encodeNilAsEmptyString {
		return Null, nil
	}
	return json.Marshal(s.String())
}

func (s *String) UnmarshalJSON(data []byte) error {
	s.value = nil
	b := RawJSON(data)
	switch {
	case b.IsNull():
	case b.IsString():
		str := b.String()
		s.value = &str
	default:
		// TODO: really need to do better with errors
		return &json.UnmarshalTypeError{}
	}
	return nil
}

func formatString(value interface{}, layout ...string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case string:
		return &v, nil
	case []byte:
		return formatString(string(v))
	case String:
		return formatString(v.String())
	case *String:
		return formatString(*v)
	case time.Time:
		if len(layout) > 0 {
			return formatString(v.Format(layout[0]))
		}
		return formatString(v.String())
	case *time.Time:
		return formatString(*v)
	case fmt.Stringer:
		return formatString(v.String())
	case bool:
		return formatString(strconv.FormatBool(v))
	case int64:
		return formatString(strconv.FormatInt(v, 10))
	case int32:
		return formatString(int64(v))
	case int16:
		return formatString(int64(v))
	case int8:
		return formatString(int64(v))
	case int:
		return formatString(int64(v))
	case uint64:
		return formatString(strconv.FormatUint(v, 10))
	case uint32:
		return formatString(uint64(v))
	case uint16:
		return formatString(uint64(v))
	case uint8:
		return formatString(uint64(v))
	case uint:
		return formatString(uint64(v))
	case float64:
		return formatString(strconv.FormatFloat(v, 'f', 0, 64))
	case float32:
		return formatString(strconv.FormatFloat(float64(v), 'f', 0, 32))
	case complex128:
		return formatString(strconv.FormatComplex(v, 'f', 0, 128))
	case complex64:
		return formatString(strconv.FormatComplex(complex128(v), 'f', 0, 64))
	case []string:
		return formatString(strings.Join(v, ","))
	case *string:
		return formatString(*v)
	case *bool:
		return formatString(*v)
	case *int64:
		return formatString(*v)
	case *int32:
		return formatString(*v)
	case *int16:
		return formatString(*v)
	case *int8:
		return formatString(*v)
	case *int:
		return formatString(*v)
	case *uint64:
		return formatString(*v)
	case *uint32:
		return formatString(*v)
	case *uint16:
		return formatString(*v)
	case *uint8:
		return formatString(*v)
	case *uint:
		return formatString(*v)
	case *float64:
		return formatString(*v)
	case *float32:
		return formatString(*v)
	case *complex128:
		return formatString(*v)
	case *complex64:
		return formatString(*v)
	default:
		return nil, ErrInvalidType
	}
}
func (s *String) Int() (int64, error) {
	str := s.String()
	return strconv.ParseInt(str, 0, 64)
}
func (s *String) Uint() (uint64, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return strconv.ParseUint(str, 0, 64)
}

func (s *String) Bytes() []byte {
	if s == nil {
		return nil
	}
	if s.IsNil() {
		return nil
	}
	return []byte(*s.value)
}

func (s *String) Float() (float64, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return strconv.ParseFloat(str, 64)
}
func (s *String) Float64() (float64, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return strconv.ParseFloat(str, 64)
}
func (s *String) Float32() (float32, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	v, err := strconv.ParseFloat(str, 32)
	return float32(v), err
}

func (s *String) Time(layout string) (time.Time, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return time.Parse(layout, str)
}

func (s *String) Duration() (time.Duration, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return time.ParseDuration(str)
}

func (s *String) Complex() (complex128, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return strconv.ParseComplex(str, 128)
}
func (s *String) Complex128() (complex128, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return strconv.ParseComplex(str, 128)
}
func (s *String) Complex64() (complex64, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	v, err := strconv.ParseComplex(str, 64)
	return complex64(v), err
}

func (s *String) Bool() (bool, error) {
	var str string
	if s != nil {
		str = s.String()
	}
	return strconv.ParseBool(str)
}

func (s *String) Set(value interface{}, timeLayout ...string) error {
	s.value = nil
	if value == nil {
		return nil
	}
	v, err := formatString(value, timeLayout...)
	if err != nil {
		return err
	}
	s.value = v
	return nil
}
