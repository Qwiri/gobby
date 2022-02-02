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
