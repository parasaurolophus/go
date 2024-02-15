// Copyright Kirk Rader 2024

package stacktraces

import (
	"fmt"
	"runtime"
)

// Return the name of the function that called this one, i.e. the currently
// executing function from that function's point of view.
func FunctionName() string {

	name, _, _, _ := functionInfo(nil)
	return name
}

// Return the name of the function at the specified position in the current call
// stack and nil, or the empty string and a StackTrack if no such function is
// found.
func FunctionNameAt(skipFrames any) (string, error) {

	name, _, _, err := functionInfo(skipFrames)
	return name, err
}

// Return name, source file name and line number of the function that this one.
func FunctionInfo(skipFrames any) (string, string, int, error) {

	return functionInfo(skipFrames)
}

// Common implementation for FunctionName() and FunctionInfo()
func functionInfo(skipFrames any) (string, string, int, error) {

	var (
		skip                     = 0
		frameTest stackFrameTest = func(frame *runtime.Frame) bool { return true }
	)

	switch v := skipFrames.(type) {

	case int:
		if v < 0 {
			skip = defaultSkip
		} else {
			skip = v
		}

	case string:
		frameTest = skipUntil(v)

	default:
		skip = defaultSkip
	}

	pc := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	for {
		frame, more := frames.Next()

		if frameTest(&frame) {
			return frame.Function, frame.File, frame.Line, nil
		}

		if !more {
			return "", "", 0, New(fmt.Sprintf("no frame found for %v", skipFrames), nil)
		}
	}
}
