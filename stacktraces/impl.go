package stacktraces

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
)

type (

	// Type of function used to write a single frame to a stack trace.
	stackFrameFormatter func(writer *bufio.Writer, frameNumber int, frame *runtime.Frame)

	// Type of function used to test a frame for inclusion in a stack trace.
	stackFrameTest func(frame *runtime.Frame) bool
)

// Return a function which writes a multi-line representation of a stack frame.
func longFrameFormatter() stackFrameFormatter {

	writeSeparator := false

	return func(writer *bufio.Writer, frameNumber int, frame *runtime.Frame) {

		if writeSeparator {
			writer.WriteString("---\n")
		} else {
			writeSeparator = true
		}

		writer.WriteString(
			fmt.Sprintf(
				"%d:%s\n%s:%d\n%#v\n",
				frameNumber,
				frame.Function,
				frame.File,
				frame.Line,
				frame.PC))
	}
}

// Return a function which writes a one-line representation of a stack frame.
func shortFrameFormatter() stackFrameFormatter {

	writeSeparator := false

	return func(writer *bufio.Writer, frameNumber int, frame *runtime.Frame) {

		if writeSeparator {
			writer.WriteString(" < ")
		} else {
			writeSeparator = true
		}

		writer.WriteString(
			fmt.Sprintf(
				"%d:%s [%s:%d]",
				frameNumber,
				frame.Function,
				frame.File,
				frame.Line))
	}
}

// Helper used by New(string, any), LongStackTrace(any) and
// ShortStackTrace(any).
//
// Returns long and short string representations of the current call stack.
//
//   - skipFrames may be an int or string.
//
//   - startAfter is a function name that is used if and only if skipFrames is
//     a negative int.
//
//   - frameWriter is a function that is used to append stack frames to the
//     output string.
//
// This function operates by:
//
//  1. Creating a test function based on the combination of values for skipFrames
//     and startAfter.
//
//  2. Iterating over the call stack and calling frameWriter for each frame for
//     which the test function returns true.
//
// The combination of skipFrames and startAfter are used to select both a
// numerical number of frames to skip before the iteration begins as well as the
// name of the function to use as the "skip until a matching frame is found"
// logic is specified. When skipFrames is a string, the test function is created
// using startWhenSeen(string), passing the value of skipFrames as the paramter.
// When skipFrames is a negative int, the test function is created using
// startAfterSeen(string), passing the value of startAfter as a parameter.
// Otherwise, the test function defaults to one which unconditionally returns
// true and the numeric value of skipFrames is passed to runtime.Callers as its
// "skip" parameter.
//
// See startWhenSeen(string) for the implementations of the test function used
// when skipFrames is a string.
func formatStackTrace(skipFrames any, longFormatter stackFrameFormatter, shortFormatter stackFrameFormatter) (string, string) {

	longBuffer := bytes.Buffer{}
	longWriter := bufio.NewWriter(&longBuffer)

	shortBuffer := bytes.Buffer{}
	shortWriter := bufio.NewWriter(&shortBuffer)

	skip := 0
	frameTest := func(*runtime.Frame) bool { return true }

	switch v := skipFrames.(type) {

	case int:
		if v < 0 {
			// the number of frames to skip is empirically derived and may
			// change any time the code in this library is refactored
			skip = 3
		} else {
			skip = v
		}

	case string:
		frameTest = startWhenSeen(v)
	}

	pc := make([]uintptr, 1024)
	n := runtime.Callers(skip, pc)

	if n < 1 {
		return "", ""
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	frameNumber := skip

	for {

		frame, more := frames.Next()

		if frameTest(&frame) {

			if longFormatter != nil {
				longFormatter(longWriter, frameNumber, &frame)
			}

			if shortFormatter != nil {
				shortFormatter(shortWriter, frameNumber, &frame)
			}
		}

		frameNumber += 1

		if !more {
			break
		}
	}

	longWriter.Flush()
	shortWriter.Flush()
	return longBuffer.String(), shortBuffer.String()
}

// Return a function that returns true when invoked for a frame the given
// function name and all that follow it.
//
// This is the frame test used by formatStackTrace() when skipFrames is a string
func startWhenSeen(startWhen string) stackFrameTest {

	seen := false

	return func(frame *runtime.Frame) bool {

		seen = seen || startWhen == frame.Function
		return seen
	}
}
