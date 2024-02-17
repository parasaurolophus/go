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
	StackTrace struct {
		// String returned by StackTrace.Error().
		msg string
		// Multi-line stack trace returned by StackTrace.LongTrace().
		longTrace string
		// One-line stack trace returned by StackTrace.ShortTrace().
		shortTrace string
	}
	// Type of function used to write a single frame to a stack trace.
	stackFrameFormatter func(writer *bufio.Writer, frameNumber int, frame *runtime.Frame)
	// Type of function used to test a frame for inclusion in a stack trace.
	stackFrameTest func(frame *runtime.Frame) bool
)

const (

	// Capacity of address buffer passed to runtime.Callers().
	maxDepth = 1024
	// The default number of frames to skip when skipFrames is neither a string
	// nor a positive int.
	//
	// This could change as a result of refactoring this or Go's standard
	// runtime libraries. Current value of 3 is based on:
	//
	// 1. The function where the trace should start calls FunctionInfo(),
	//    FunctionNameAt(), LongStackTrace(), ShortStackTrace() or New().
	//
	// 2. LongstackTrace(), ShortStackTrace() and New() call formatStackTrace()
	//
	// 3. formatStackTrace() calls runtime.Callers() and that is the first frame
	//    on the stack.
	defaultSkip = 3
)

// Convenience function for creating a multi-line trace of the current call
// stack.
//
// See New() for a description of the skipFrames parameter.
func LongStackTrace(skipFrames any) string {
	trace, _ := formatStackTrace(skipFrames, longFrameFormatter(), nil)
	return trace
}

// Convenience function for creating a one-line trace of the current call stack.
//
// See New() for a description of the skipFrames parameter.
func ShortStackTrace(skipFrames any) string {
	_, trace := formatStackTrace(skipFrames, nil, shortFrameFormatter())
	return trace
}

// Return a newly created StackTrace object that captures the current call
// stack.
//
// The given msg will be returned by StackTrace.Error().
//
// skipFrames may be an int or string; when passed a value of any other type,
// the actual value will be ignored and the behavior will be as described below.
//
// The strings returned by stackTrace.LongTrace() and stackTrace.ShortTrace()
// will include all of the call stack frames at the time this constructor is
// called, excluding ones from the top of the stack as follows:
//
//   - If skipFrames is a non-negative int the specified number of frames are
//     skipped.
//
//   - If skipFrames is a string all frames before the function with the
//     given name are skipped.
//
//   - If skipFrames is any other value, all frames up to and including this
//     function's frame are skipped.
//
// The empty string is returned for both traces if the stack depth is exceeded
// when passing a positive int or no matching frame is found when passing a
// string.
func New(msg string, skipFrames any) StackTrace {
	longFormatter := longFrameFormatter()
	shortFormatter := shortFrameFormatter()
	long, short := formatStackTrace(skipFrames, longFormatter, shortFormatter)
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
// Returns multi-line and one-line string representations of the current call
// stack.
func formatStackTrace(skipFrames any, longFormatter stackFrameFormatter, shortFormatter stackFrameFormatter) (string, string) {
	// in-memory writers in which the stack trace strings will be built
	longBuffer := bytes.Buffer{}
	longWriter := bufio.NewWriter(&longBuffer)
	shortBuffer := bytes.Buffer{}
	shortWriter := bufio.NewWriter(&shortBuffer)
	// the number that will be passed to runtime.Callers()
	skip := 0
	// default frame test unconditionally returns true on the assumption that
	// skip will be a non-negative int
	frameTest := func(*runtime.Frame) bool { return true }
	// adjust skip and frameTest according to the supplied value of skipFrames
	switch v := skipFrames.(type) {
	case int:
		if v < 0 {
			// skip past this function's caller's caller when skipFrames is
			// negative
			skip = defaultSkip
		} else {
			// skip the specified number of frames when skipFrames is
			// non-negative
			skip = v
		}
	case string:
		// use a frameTest that skips all frames until a given function is seen
		frameTest = skipUntil(v)
	default:
		// skip past this function's caller's caller when skipFrames is any
		// other value
		skip = defaultSkip
	}
	// capture the current call stack's addresses using (what is to be hoped) a
	// sufficient size of address buffer
	pc := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pc)
	// when stack depth was not exceeded...
	if n > 0 {
		// ...iterate over the captured call stack's frames, applying the
		// formatters to each frame for which frameTest returns true
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
	}
	longWriter.Flush()
	shortWriter.Flush()
	return longBuffer.String(), shortBuffer.String()
}

// Return a function that returns true when invoked for a frame with the given
// function name and all that follow it.
//
// The returned function will be the frame test used by formatStackTrace() when
// skipFrames is a string
func skipUntil(startWhen string) stackFrameTest {
	seen := false
	return func(frame *runtime.Frame) bool {
		seen = seen || startWhen == frame.Function
		return seen
	}
}
