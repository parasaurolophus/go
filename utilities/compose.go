// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

// Invoke the given function asynchronously for each value sent to the given
// input channel, sending the result to the given output channel.
func Async[T any](f func(T) T, in chan T, out chan T) {
	defer close(out)
	options := logging.LoggerOptions{
		BaseTags: []string{"utilities.Async", logging.PANIC},
	}
	logger := logging.New(os.Stderr, &options)
	invoke := func(arg T) {
		defer func() {
			if r := recover(); r != nil {
				logger.Always(
					func() string {
						return fmt.Sprintf("recovered from panic in injected dependency: %v", r)
					})
			}
		}()
		out <- f(arg)
	}
	for arg := range in {
		invoke(arg)
	}
}

// Return the result of invoking the given list of functions. The first function
// is applied to the given value. The second to the result returned by the
// first, and so on. The final value and nil is returned by Compose if all
// functions execute successfully. The last successful value and an error are
// returned for the first function that panics.
func Compose[T any](value T, functions ...func(T) T) (T, error) {
	var err error
	invoke := func(f func(T) T, p T) T {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v: %s", r, stacktraces.ShortStackTrace(nil))
			}
		}()
		return f(p)
	}
	for _, fn := range functions {
		v := invoke(fn, value)
		if err != nil {
			return value, err
		}
		value = v
	}
	return value, err
}

// Return a new slice containing the results of applying the specified function
// to each element of the given one.
func Map[T any](f func(T) T, slice []T) []T {

	result := []T{}

	for _, element := range slice {
		result = append(result, f(element))
	}

	return result
}
