_Copyright &copy; Kirk Rader 2024_

# Stack Traces for Logging and Error Messages in Go

Note that the primary documentation for this package is inline comments in the
source code, viewable by way of `go doc`.

```bash
cd stacktraces
go doc -all
```

```
package stacktraces // import "parasaurolophus/go/stacktraces"


FUNCTIONS

func FunctionInfo(skipFrames any) (string, string, int, error)
    Return name, source file name and line number of the function that this one.

func FunctionName() string
    Return the name of the function that called this one, i.e. the currently
    executing function from that function's point of view.

func FunctionNameAt(skipFrames any) (string, error)
    Return the name of the function at the specified position in the current
    call stack and nil, or the empty string and a StackTrack if no such function
    is found.

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

## Examples

```go
// stacktraces.StackTrace is an error object that can be used in place of
// errors.New(message)
err := func() error {
  return stacktraces.New("message", nil)
}()

fmt.Printf("err.Error(): %s\n\n", err.Error())

// in addition, stacktraces.StackTrace captures one-line and multi-line
// stack traces
stackTrace := err.(stacktraces.StackTrace)
fmt.Printf("stackTrace.ShortTrace(): %s\n\n", stackTrace.ShortTrace())
fmt.Printf("stackTrace.LongTrace():\n%s\n", stackTrace.LongTrace())

// stacktraces.LongStackTrace() and stacktraces.ShortStackTrace() helper
// functions are also provided
fmt.Printf("stacktraces.ShortStackTrace(0): %s\n\n", stacktraces.ShortStackTrace(0))
fmt.Printf("stacktraces.LongStackTrace(0):\n%s", stacktraces.LongStackTrace(0))
```

writes the following to `stdout`:

```
err.Error(): message

stackTrace.ShortTrace(): 3:main.main.func1 [/source/go/scratch/scratch.go:15] < 4:main.main [/source/go/scratch/scratch.go:16] < 5:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 6:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]

stackTrace.LongTrace():
3:main.main.func1
/source/go/scratch/scratch.go:15
0x8e38b
---
4:main.main
/source/go/scratch/scratch.go:16
0x8e370
---
5:runtime.main
/usr/local/go/src/runtime/proc.go:267
0x431ab
---
6:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1197
0x6dc73

stacktraces.ShortStackTrace(0): 0:runtime.Callers [/usr/local/go/src/runtime/extern.go:308] < 1:parasaurolophus/go/stacktraces.formatStackTrace [/source/go/stacktraces/stacktraces.go:250] < 2:parasaurolophus/go/stacktraces.ShortStackTrace [/source/go/stacktraces/stacktraces.go:73] < 3:main.main [/source/go/scratch/scratch.go:28] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]

stacktraces.LongStackTrace(0):
0:runtime.Callers
/usr/local/go/src/runtime/extern.go:308
0x8dfd3
---
1:parasaurolophus/go/stacktraces.formatStackTrace
/source/go/stacktraces/stacktraces.go:250
0x8df80
---
2:parasaurolophus/go/stacktraces.LongStackTrace
/source/go/stacktraces/stacktraces.go:64
0x8d637
---
3:main.main
/source/go/scratch/scratch.go:29
0x8e5a3
---
4:runtime.main
/usr/local/go/src/runtime/proc.go:267
0x431ab
---
5:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1197
0x6dc73
```
