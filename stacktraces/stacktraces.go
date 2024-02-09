// Copyright Kirk Rader 2024

package stacktraces

import (
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
