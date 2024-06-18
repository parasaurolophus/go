// Copyright Kirk Rader 2024

package stacktraces_test

import (
	"errors"
	"strconv"
	"strings"
)

// Return the function name and frame number from the first entry in the given
// output of stacktraces.ShortStackTrace(any).
func FirstFunctionShort(strackTrace string) (string, int, error) {

	frame := strings.Split(strackTrace, "<")[0]
	parts := strings.Split(frame, ":")

	if len(parts) < 2 {
		return "", 0, errors.New("no colon in stack frame")
	}

	name := strings.Trim(strings.Split(parts[1], "[")[0], " ")
	n, err := strconv.Atoi(parts[0])
	return name, n, err
}

// Return the function name and frame number from the first entry in the given
// output of stacktraces.LongStackTrace(any).
func FirstFunctionLong(stackTrace string) (string, int, error) {

	frame := strings.Split(stackTrace, "\n")[0]
	parts := strings.Split(frame, ":")

	if len(parts) < 2 {
		return "", 0, errors.New("no colon in stack frame")
	}

	name := parts[1]
	n, err := strconv.Atoi(parts[0])
	return name, n, err
}
