package calling

import (
	"fmt"
	"reflect"
)

var ErrorType = reflect.TypeOf((*error)(nil)).Elem()

type WithStringArgsFunc func(args ...string)
type WithStringArgsErrorFunc func(args ...string) error

func WithStringArgs(function interface{}) WithStringArgsFunc {
	v := reflect.ValueOf(function)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic("not a function")
	}
	if t.NumOut() != 0 {
		panic("must not return results")
	}
	numArgs := t.NumIn()
	argTypes := make([]reflect.Type, numArgs)
	for i := range argTypes {
		argTypes[i] = t.In(i)
	}
	return func(stringArgs ...string) {
		if len(stringArgs) != numArgs {
			panic("number of string args is no equal number of target function args")
		}
		args := make([]reflect.Value, numArgs)
		for i := range args {
			args[i] = reflect.Zero(argTypes[i])
			_, err := fmt.Sscan(stringArgs[i], args[i].Interface())
			if err != nil {
				panic(fmt.Errorf("Could not convert string argument %d '%s' to type %s becuase of error: %s", i, stringArgs[i], argTypes[i], err))
			}
		}
		v.Call(args)
	}
}

func WithStringArgsError(function interface{}) WithStringArgsErrorFunc {
	v := reflect.ValueOf(function)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic("not a function")
	}
	if t.NumOut() != 1 || t.Out(0) != ErrorType {
		panic("must return an error")
	}
	numArgs := t.NumIn()
	argTypes := make([]reflect.Type, numArgs)
	for i := range argTypes {
		argTypes[i] = t.In(i)
	}
	return func(stringArgs ...string) error {
		if len(stringArgs) != numArgs {
			panic("number of string args is no equal number of target function args")
		}
		args := make([]reflect.Value, numArgs)
		for i := range args {
			args[i] = reflect.Zero(argTypes[i])
			_, err := fmt.Sscan(stringArgs[i], args[i].Interface())
			if err != nil {
				panic(fmt.Errorf("Could not convert string argument %d '%s' to type %s becuase of error: %s", i, stringArgs[i], argTypes[i], err))
			}
		}
		return v.Call(args)[0].Interface().(error)
	}
}
