package validate

import "errors"

var (
	ErrNumberNotNumber = errors.New("number is not a number")
	ErrNumberTooSmall  = errors.New("number is too small")
	ErrNumberTooLarge  = errors.New("number is too large")
)

type NumberScheme struct {
	name string
	min  *int64
	max  *int64
}

func Number() *NumberScheme {
	return &NumberScheme{}
}

func (ns *NumberScheme) Name() string {
	return ns.name
}

func (ns *NumberScheme) As(name string) *NumberScheme {
	ns.name = name
	return ns
}

func (ns *NumberScheme) Min(min int64) *NumberScheme {
	ns.min = &min
	return ns
}

func (ns *NumberScheme) Max(max int64) *NumberScheme {
	ns.max = &max
	return ns
}

func (ns *NumberScheme) Exact(exact int64) *NumberScheme {
	return ns.RangeClosed(exact, exact)
}

func (ns *NumberScheme) Range(min, max int64) *NumberScheme {
	ns.min = &min
	max--
	ns.max = &max
	return ns
}

func (ns *NumberScheme) RangeClosed(min, max int64) *NumberScheme {
	ns.min, ns.max = &min, &max
	return ns
}

func (ns *NumberScheme) Check(i interface{}) error {
	var num int64
	// is this redundant?
	switch t := i.(type) {
	case int:
		num = int64(t)
	case uint:
		num = int64(t)
	case int8:
		num = int64(t)
	case uint8:
		num = int64(t)
	case int16:
		num = int64(t)
	case uint16:
		num = int64(t)
	case int32:
		num = int64(t)
	case uint32:
		num = int64(t)
	case uint64:
		num = int64(t)
	default:
		return ErrNumberNotNumber
	}
	// check min
	if ns.min != nil && num < *ns.min {
		return ErrNumberTooSmall
	}
	if ns.max != nil && num > *ns.max {
		return ErrNumberTooLarge
	}
	return nil
}
