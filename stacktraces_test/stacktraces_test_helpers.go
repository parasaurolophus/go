// Copyright Kirk Rader 2024

package stacktraces_test

import (
	"errors"
	"strconv"
	"strings"
)

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
