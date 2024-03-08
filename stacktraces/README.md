_Copyright &copy; Kirk Rader 2024_

# Stack Traces for Logging and Error Messages in Go

```
$ go doc -all
package stacktraces // import "parasaurolophus/go/stacktraces"


FUNCTIONS

func FunctionInfo(skipFrames any) (int, string, string, int, bool)
    Return call stack info for a function that directly or indirectly called
    this one.

    See README.md for an explanation of the skipFrames parameter.

    Return values are:

      - Frame number of the caller
      - Function name of the caller
      - Source file name of the caller
      - Source line number of the caller
      - true or false depending on whether or not the specified caller was
        actually found

    If the last value is false, the others are set to their "zero values."

func FunctionName() string
    Return the name of the function that called this one, i.e. the currently
    executing function from that function's point of view.

func LongStackTrace(skipFrames any) string
    Convenience function for creating a multi-line trace of the current call
    stack.

    See New() for a description of the skipFrames parameter.

func ShortStackTrace(skipFrames any) string
    Convenience function for creating a one-line trace of the current call
    stack.

    See New() for a description of the skipFrames parameter.


TYPES

type StackTrace struct {
	// Has unexported fields.
}
    Error objects that contain the current call stack at the time they were
    created.

    Use stacktraces.New(msg string, skipFrames any) in place of errors.new(msg
    string).

func New(msg string, skipFrames any) StackTrace
    Return a newly created StackTrace object that captures the current call
    stack.

    The given msg will be returned by StackTrace.Error().

    skipFrames may be an int or string; when passed a value of any other type,
    the actual value will be ignored and the behavior will be as described
    below.

    The strings returned by stackTrace.LongTrace() and stackTrace.ShortTrace()
    will include all of the call stack frames at the time this constructor is
    called, excluding ones from the top of the stack as follows:

      - If skipFrames is a non-negative int the specified number of frames are
        skipped.

      - If skipFrames is a string all frames before the function with the given
        name are skipped.

      - If skipFrames is any other value, all frames up to and including this
        function's frame are skipped.

    The empty string is returned for both traces if the stack depth is exceeded
    when passing a positive int or no matching frame is found when passing a
    string.

func (t StackTrace) Error() string
    Implement error interface.

func (t StackTrace) LongTrace() string
    Return the multi-line representation of this stack trace.

func (t StackTrace) ShortTrace() string
    Return the one-line representation of this stack trace.
```
