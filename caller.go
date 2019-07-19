package gosocketio

import (
	"errors"
	"fmt"
	"reflect"
)

type caller struct {
	Func        reflect.Value
	Args        []reflect.Type
	ArgsPresent bool
	Out         bool
}

var (
	ErrorCallerNotFunc     = errors.New("f is not function")
	ErrorCallerNot2Args    = errors.New("f should have 1 or 2 args")
	ErrorCallerMaxOneValue = errors.New("f should return not more than one value")
)

/**
Parses function passed by using reflection, and stores its representation
for further call on message or ack
*/
func newCaller(f interface{}) (*caller, error) {
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		return nil, ErrorCallerNotFunc
	}

	fType := fVal.Type()
	if fType.NumOut() > 1 {
		return nil, ErrorCallerMaxOneValue
	}

	curCaller := &caller{
		Func: fVal,
		Out:  fType.NumOut() == 1,
	}
	if fType.NumIn() == 1 {
		curCaller.Args = []reflect.Type{}
		curCaller.ArgsPresent = false
	} else if fType.NumIn() >= 2 {
		curCaller.Args = make([]reflect.Type, fType.NumIn()-1)
		for i := 1; i < fType.NumIn(); i++ {
			curCaller.Args[i-1] = fType.In(i)
		}
		curCaller.ArgsPresent = true
	} else {
		return nil, ErrorCallerNot2Args
	}

	return curCaller, nil
}

/**
returns function parameter as it is present in it using reflection
*/
func (c *caller) getArgs() []interface{} {
	r := make([]interface{}, len(c.Args))
	for index := range c.Args {
		r[index] = reflect.New(c.Args[index]).Interface()
	}
	return r
}

/**
calls function with given arguments from its representation using reflection
*/
func (c *caller) callFunc(h *Channel, args ...interface{}) []reflect.Value {
	//nil is untyped, so use the default empty value of correct type
	if args == nil {
		args = c.getArgs()
	}

	a := make([]reflect.Value, len(args)+1)
	a[0] = reflect.ValueOf(h)
	for i := range args {
		fmt.Println(args[i])
		if args[i] != nil {
			a[i+1] = reflect.ValueOf(args[i]).Elem()
		} else {
			// a[i+1] =
		}
	}
	// a := []reflect.Value{reflect.ValueOf(h), reflect.ValueOf(args).Elem()}
	if !c.ArgsPresent {
		a = a[0:1]
	}
	fmt.Println("Call Argument size", a, c.Func)

	return c.Func.Call(a)
}
