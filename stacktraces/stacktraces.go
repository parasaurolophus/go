// Copyright Kirk Rader 2024

package stacktraces

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
)

type (
	// Error objects that contain the current call stack at the time they were
	// created.
	//
	// Use stacktraces.New(msg string, skipFrames any) in place of
	// errors.new(msg string).
	//
	// See LongStackTrace(skipFrames any) and ShortStackTrace(skipFrames any).
	StackTrace struct {

		// String returned by StackTrace.Error().
		msg string

		// String returned by StackTrace.LongTrace().
		longTrace string

		// String returned by StackTrace.ShortTrace().
		shortTrace string
	}
)

// Return the name of the function that called this one, i.e. the currently
// executing function from that function's point of view.
func FunctionName() string {

	// the number of frames to skip is empirically derived and may change as a
	// result of any refactoring of this library

	// current assumes first frame is runtime.Callers(), second frame is this
	// function

	pc := make([]uintptr, 1024)
	n := runtime.Callers(2, pc)
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return frame.Function
}

// Return a newly created object that implements the StackTrace interface.
//
// The given msg will be returned by StackTrace.Error().
//
// skipFrames may be an int or string; when passed a value of any other type,
// the value will be ignored and the behavior will be as if 0 were passed.
//
// The traces returned by stackTrace.LongTrace() and stackTrace.ShortTrace()
// will include all of the call stack frames at the time this constructor is
// called, excluding ones from the top of the stack as follows:
//
//   - If skipFrames is a non-negative int the specified number of frames are
//     skipped.
//
//   - If skipFrames is a string all frames before the frame for the function
//     with the given name are skipped.
//
//   - If skipFrames is any other value, all frames up to and including this
//     function's frame are skipped.
//
// The empty string is returned if the stack depth is exceeded when passing a
// positive int or no matching frame is found when passing a string.
func New(msg string, skipFrames any) StackTrace {

	longFormatter := longFrameFormatter()
	shortFormatter := shortFrameFormatter()
	long := ""
	short := ""

	switch v := skipFrames.(type) {

	case int:

		if v < 0 {

			// the number of frames to skip is empirically derived and may
			// change as a result of any refactoring of this library's code
			long, short = formatStackTrace(3, longFormatter, shortFormatter)

		} else {

			long, short = formatStackTrace(v, longFormatter, shortFormatter)
		}

	case string:

		long, short = formatStackTrace(v, longFormatter, shortFormatter)

	default:

		// the number of frames to skip is empirically derived and may
		// change as a result of any refactoring of this library's code
		long, short = formatStackTrace(3, longFormatter, shortFormatter)
	}

	return StackTrace{
		msg:        msg,
		longTrace:  long,
		shortTrace: short,
	}
}

// Implement error interface.
func (t StackTrace) Error() string {

	return t.msg
}

// Return the multi-line representation of this stack trace.
func (t StackTrace) LongTrace() string {

	return t.longTrace
}

// Return the one-line representation of this stack trace.
func (t StackTrace) ShortTrace() string {

	return t.shortTrace
}

func LongStackTrace(skipFrames any) string {

	trace, _ := formatStackTrace(skipFrames, longFrameFormatter(), nil)
	return trace
}

func ShortStackTrace(skipFrames any) string {

	_, trace := formatStackTrace(skipFrames, nil, shortFrameFormatter())
	return trace
}

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

	default:
		// the number of frames to skip is empirically derived and may
		// change any time the code in this library is refactored
		skip = 3
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

// Return a function that returns true when invoked for a frame with the given
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
