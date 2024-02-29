// Copyright Kirk Rader 2024

package closures

import (
	"fmt"
	"os"
)

type (

	// Signature of the first return value from MakeClosures().
	//
	// This will return the current value of the persistent binding in the
	// closed environment each time it is called.
	Getter func() int

	// Signature of the second return value from MakeClosures().
	//
	// This will update and then return the value of persistent binding in the
	// closed environment each time it is called.
	Setter func(newValue int) int
)

// Return two closures which encapsulate a binding to the given value in their
// shared environment.
func MakeClosures(value int) (Getter, Setter) {

	getter := func() int { return value }
	setter := func(newValue int) int {
		value = newValue
		return value
	}

	return getter, setter
}

// Demonstrate functional composition.
//
// This is an example of how to implement CPS (Continuation Passing Style) in a
// language like Go which lacks support for first-class continuations and tail
// call optimization. Each function in the chain of transformations is
// effectively the continuation of the computation from the preceding one.
//
// Note that we use the terminology of functional programming and its underlying
// mathematical model here, but this is a perfect example of what design pattern
// theorists refer to as inversion of control.
func PassContinuations[T any](value T, closures ...func(T) T) T {

	// use this with defer to log panics but then continue execution
	panicHandler := func() {

		// Invoking this function using defer guarantees that it will always be
		// called when the calling function is exiting, either normally or due
		// to a panic.

		// Calling recover() produces the value passed to panic() if the calling
		// function is exiting abnormally. The panic-induced stack-unwinding
		// then stops unless this function were to call panic() again.

		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "recovered - %v\n", r)
		}
	}

	// apply c to value with a deferred function that logs panics and recovers from
	// them
	transformValue := func(c func(T) T) {

		defer panicHandler()
		value = c(value)
	}

	for _, c := range closures {

		// wrap invocation of each passed-in closure in a panic recovery handler
		// so that one misbhaving injected dependency doesn't crash the whole
		// program
		transformValue(c)
	}

	return value
}
