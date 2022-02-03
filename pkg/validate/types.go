package validate

import (
	"errors"
	"fmt"
)

type Checker interface {
	Check(interface{}) error
}

type scheme struct {
	name string
}

func (scheme) Check(interface{}) error {
	fmt.Println("WARN | Checking default scheme")
	return nil
}

type Schemes []Checker

type Result map[string]interface{}

func (r Result) Get(name string, def interface{}) interface{} {
	val, ok := r[name]
	if !ok {
		return def
	}
	return val
}

func (r Result) String(name string) string {
	val, ok := r.Get(name, "").(string)
	if !ok {
		return ""
	}
	return val
}

func (r Result) Number(name string) int64 {
	val, ok := r.Get(name, 0).(int64)
	if !ok {
		return 0
	}
	return val
}

var (
	ErrArgLength = errors.New("argument length does not match expected scheme")
	ErrNotScheme = errors.New("invalid scheme")
)

func (ss Schemes) Check(args ...interface{}) (res Result, err error) {
	if len(args) != len(ss) {
		return nil, ErrArgLength
	}
	res = make(Result)
	for i, chk := range ss {
		sch, ok := chk.(scheme)
		if !ok {
			return nil, ErrNotScheme
		}
		arg := args[i]
		if err = chk.Check(arg); err != nil {
			return
		}
		res[sch.name] = arg
	}
	return
}
