// Package calling provides utilities for calling functions with string arguments
// that are automatically converted to the appropriate types using reflection.
//
// This package is useful for CLI applications, configuration systems, or any
// scenario where you need to call functions with arguments provided as strings.
//
// Example usage:
//
//	func add(a, b int) { fmt.Println(a + b) }
//	wrapped := calling.WithStringArgs(add)
//	wrapped("5", "3") // Prints: 8
package calling

import (
	"fmt"
	"reflect"
)

// WithStringArgsFunc is a function type that accepts string arguments
// and internally converts them to the appropriate types before calling
// the wrapped function.
type WithStringArgsFunc func(args ...string)

// WithStringArgsErrorFunc is like WithStringArgsFunc but for functions
// that return an error.
type WithStringArgsErrorFunc func(args ...string) error

// WithStringArgs wraps a function to accept string arguments that are
// automatically converted to the function's parameter types.
//
// The wrapped function must:
//   - Be a function (not a method or other type)
//   - Return no results
//
// String arguments are converted using fmt.Sscan, which supports:
//   - Basic types: int, float, bool, string, etc.
//   - Any type that implements fmt.Scanner
//
// Example:
//
//	func greet(name string, age int) {
//	    fmt.Printf("Hello %s, you are %d years old\n", name, age)
//	}
//	wrapped := calling.WithStringArgs(greet)
//	wrapped("Alice", "30") // Calls greet("Alice", 30)
//
// Panics if:
//   - function is not a function
//   - function returns any results
//   - number of string arguments doesn't match function parameters
//   - string argument cannot be converted to the expected type
func WithStringArgs(function any) WithStringArgsFunc {
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

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// WithStringArgsError wraps a function that returns an error to accept
// string arguments that are automatically converted to the function's parameter types.
//
// The wrapped function must:
//   - Be a function (not a method or other type)
//   - Return exactly one result of type error
//
// String arguments are converted using fmt.Sscan, which supports:
//   - Basic types: int, float, bool, string, etc.
//   - Any type that implements fmt.Scanner
//
// Example:
//
//	func divide(a, b int) error {
//	    if b == 0 {
//	        return errors.New("division by zero")
//	    }
//	    fmt.Println(a / b)
//	    return nil
//	}
//	wrapped := calling.WithStringArgsError(divide)
//	err := wrapped("10", "2") // Calls divide(10, 2), returns nil
//	err = wrapped("10", "0")  // Returns error: "division by zero"
//
// Panics if:
//   - function is not a function
//   - function doesn't return exactly one error
//   - number of string arguments doesn't match function parameters
//   - string argument cannot be converted to the expected type
func WithStringArgsError(function any) WithStringArgsErrorFunc {
	v := reflect.ValueOf(function)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic("not a function")
	}
	if t.NumOut() != 1 || t.Out(0) != typeOfError {
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
		err, _ := v.Call(args)[0].Interface().(error)
		return err
	}
}
