package validate

import (
	"errors"
	"regexp"
)

var (
	ErrStringTooShort    = errors.New("string too short")
	ErrStringTooLong     = errors.New("string too long")
	ErrStringNotMatching = errors.New("string does not match pattern")
	ErrStringNotString   = errors.New("string not a string")
)

type StringScheme struct {
	name      string
	minLength int
	maxLength int
	pattern   *regexp.Regexp
}

func String() *StringScheme {
	return &StringScheme{}
}

func (st *StringScheme) Name() string {
	return st.name
}

func (st *StringScheme) As(name string) *StringScheme {
	st.name = name
	return st
}

func (st *StringScheme) Min(len int) *StringScheme {
	st.minLength = len
	return st
}

func (st *StringScheme) Max(len int) *StringScheme {
	st.maxLength = len
	return st
}

func (st *StringScheme) Exact(len int) *StringScheme {
	st.minLength, st.maxLength = len, len
	return st
}

func (st *StringScheme) Pattern(pattern string) *StringScheme {
	st.pattern = regexp.MustCompile(pattern)
	return st
}

func (st *StringScheme) Check(i interface{}) error {
	against, ok := i.(string)
	if !ok {
		return ErrStringNotString
	}
	// check min length
	if st.minLength > 0 && len(against) < st.minLength {
		return ErrStringTooShort
	}
	// check max length
	if st.maxLength > 0 && len(against) > st.maxLength {
		return ErrStringTooLong
	}
	// check pattern
	if st.pattern != nil && !st.pattern.MatchString(against) {
		return ErrStringNotMatching
	}
	return nil
}
