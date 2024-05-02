// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Invoke the given function asynchronously for each value sent to the given
// input channel, sending the result to the given output channel.
func Async[T any](f func(T) T, in chan T, out chan T) {

	defer close(out)

	invoke := func(arg T) {

		defer func() {
			if r := recover(); r != nil {
				message := map[string]any{
					"msg":       fmt.Sprintf("recovered from panic in injected dependency: %v", r),
					"time":      time.Now().Format(time.RFC3339),
					"verbosity": "ALWAYS",
					"recovered": r,
					"tags":      []string{"utilities.Async", "PANIC"},
				}
				b, _ := json.Marshal(message)
				fmt.Fprintln(os.Stderr, string(b))
			}
		}()
		out <- f(arg)
	}

	for arg := range in {
		invoke(arg)
	}
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

// Apply the first of the given functions to the specified value, then the
// second to the result from the first, and so on. Return the result of the last
// invocation.
func Reduce[T any](value T, functions ...func(T) T) T {

	for _, f := range functions {

		value = f(value)
	}

	return value
}
