// Copyright Kirk Rader 2024

package stacktraces

import "runtime"

// Return the name of the function that called this one, i.e. the currently
// executing function from that function's point of view.
func FunctionName() string {

	name, _, _ := functionInfo()
	return name
}

// Return name, source file name and line number of the function that this one.
func FunctionInfo() (string, string, int) {

	return functionInfo()
}

// Common implementation for FunctionName() and FunctionInfo()
func functionInfo() (string, string, int) {

	// The number of frames to skip is empirically derived and may change as a
	// result of any refactoring of this function or Go's standard runtime
	// library.
	//
	// The current value assumes that the first frame is runtime.Callers(), the
	// second frame is this function and the third frame is one of the public
	// wrappers for this function.
	const skip = 3

	pc := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return frame.Function, frame.File, frame.Line
}
