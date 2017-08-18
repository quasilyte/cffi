package funcall

import (
	"reflect"
)

type (
	// Func is a wrapped C function which can be
	// called with Go types directly.
	Func func(...interface{}) interface{}
	// FuncInt is like "Func", but return type differs.
	FuncInt func(...interface{}) int
	// FuncString is like "Func", but return type differs.
	FuncString func(...interface{}) string

	// User-defined function that converts Go value to C value.
	go2CFunc func(interface{}) interface{}
	// User-defined function that converts C value to Go value.
	c2GoFunc func(interface{}) interface{}
)

// Invoker exposes methods that can invoke C functions with Go types.
type Invoker struct {
	go2c go2CFunc
	c2go c2GoFunc
}

// NewInvoker creates C functions invoker that can handle some
// Go->C and C->Go types convertions.
// The set of supported conversions depends on the provided
// callbacks.
func NewInvoker(go2c go2CFunc, c2go c2GoFunc) *Invoker {
	return &Invoker{
		go2c: go2c,
		c2go: c2go,
	}
}

func (inv *Invoker) apply(fn reflect.Value, args []interface{}) interface{} {
	xs := inv.mapArgs(args)
	return inv.c2go(fn.Call(xs)[0].Interface())
}

// Apply is like "Call", but accepts a slice of arguments instead of being
// variadic over interface{}.
func (inv *Invoker) Apply(callable interface{}, args []interface{}) interface{} {
	return inv.apply(reflect.ValueOf(callable), args)
}

// Call invokes "callable" argument with provided "args".
// Each arg is expected to be a Go value which can be converted
// to C value with Invoker go2c callback.
// Return value should be convertible to Go type with Invoker c2go callback.
func (inv *Invoker) Call(callable interface{}, args ...interface{}) interface{} {
	return inv.Apply(callable, args)
}

// Int is a convenience wrapper around "Call" which does type assertion for you.
func (inv *Invoker) Int(callable interface{}, args ...interface{}) int {
	return inv.Apply(callable, args).(int)
}

// String is a convenience wrapper around "Call" which does type assertion for you.
func (inv *Invoker) String(callable interface{}, args ...interface{}) string {
	return inv.Apply(callable, args).(string)
}

// Convert Go values slice to C values slice.
func (inv *Invoker) mapArgs(args []interface{}) []reflect.Value {
	xs := make([]reflect.Value, 0, len(args))
	for i := range args {
		x := inv.go2c(args[i])
		// To support 1->N argument mapping,
		// handle slices as a tuple-like argument
		// that should be unwrapped into len(x) C args.
		// No need for recursive translation here.
		if y, ok := x.([]interface{}); ok {
			for j := range y {
				xs = append(xs, reflect.ValueOf(y[j]))
			}
		} else {
			xs = append(xs, reflect.ValueOf(x))
		}
	}
	return xs
}

// Wrap returns a closure that can be called later.
// Uses provided invoker to perform function calling
// and arguments conversion.
func Wrap(inv *Invoker, callable interface{}) Func {
	fn := reflect.ValueOf(callable)
	return func(args ...interface{}) interface{} {
		return inv.apply(fn, args)
	}
}

// WrapInt is like "Wrap", but does result type assertion to int.
func WrapInt(inv *Invoker, callable interface{}) FuncInt {
	fn := reflect.ValueOf(callable)
	return func(args ...interface{}) int {
		return inv.apply(fn, args).(int)
	}
}

// WrapString is like "Wrap", but does result type assertion to string.
func WrapString(inv *Invoker, callable interface{}) FuncString {
	fn := reflect.ValueOf(callable)
	return func(args ...interface{}) string {
		return inv.apply(fn, args).(string)
	}
}
