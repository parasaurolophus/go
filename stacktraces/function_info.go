// Copyright Kirk Rader 2024

package stacktraces

import (
	"runtime"
)

// Return the name of the function that called this one, i.e. the currently
// executing function from that function's point of view.
func FunctionName() string {
	_, name, _, _, _ := functionInfo(nil)
	return name
}

// Return call stack info for a function that directly or indirectly called this
// one.
//
// See README.md for an explanation of the skipFrames parameter.
//
// Return values are:
//
//   - Frame number of the caller
//   - Function name of the caller
//   - Source file name of the caller
//   - Source line number of the caller
//   - true or false depending on whether or not the specified caller was actually
//     found
//
// If the last value is false, the others are set to their "zero values."
func FunctionInfo(skipFrames any) (int, string, string, int, bool) {
	return functionInfo(skipFrames)
}

// Common implementation for FunctionName() and FunctionInfo()
func functionInfo(skipFrames any) (int, string, string, int, bool) {
	var (
		skip                     = 0
		frameTest stackFrameTest = func(frame *runtime.Frame) bool { return true }
	)
	switch v := skipFrames.(type) {
	case int:
		if v < 0 {
			skip = defaultSkip - v
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
	if n < 1 {
		return 0, "", "", 0, false
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	count := skip
	for {
		frame, more := frames.Next()
		if frameTest(&frame) {
			return count, frame.Function, frame.File, frame.Line, true
		}
		if !more {
			return 0, "", "", 0, false
		}
		count += 1
	}
}
