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

func FunctionName() string
    Return the name of the function that called this one, i.e. the currently
    executing function from that function's point of view.

func LongStackTrace(skipFrames any) string
func ShortStackTrace(skipFrames any) string

TYPES

type StackTrace struct {
	// Has unexported fields.
}
    Error objects that contain the current call stack at the time they were
    created.

    Use stacktraces.New(msg string, skipFrames any) in place of errors.new(msg
    string).

    See LongStackTrace(skipFrames any) and ShortStackTrace(skipFrames any).

func New(msg string, skipFrames any) StackTrace
    Return a newly created object that implements the StackTrace interface.

    The given msg will be returned by StackTrace.Error().

    skipFrames may be an int or string; when passed a value of any other type,
    the value will be ignored and the behavior will be as if 0 were passed.

    The traces returned by stackTrace.LongTrace() and stackTrace.ShortTrace()
    will include all of the call stack frames at the time this constructor is
    called, excluding ones from the top of the stack as follows:

      - If skipFrames is a non-negative int the specified number of frames are
        skipped.

      - If skipFrames is a string all frames before the frame for the function
        with the given name are skipped.

      - If skipFrames is any other value, all frames up to and including this
        function's frame are skipped.

    The empty string is returned if the stack depth is exceeded when passing a
    positive int or no matching frame is found when passing a string.

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
0x8e47b
---
4:main.main
/source/go/scratch/scratch.go:16
0x8e460
---
5:runtime.main
/usr/local/go/src/runtime/proc.go:267
0x431ab
---
6:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1197
0x6dc73

stacktraces.ShortStackTrace(0): 0:runtime.Callers [/usr/local/go/src/runtime/extern.go:308] < 1:parasaurolophus/go/stacktraces.formatStackTrace [/source/go/stacktraces/stacktraces.go:235] < 2:parasaurolophus/go/stacktraces.ShortStackTrace [/source/go/stacktraces/stacktraces.go:139] < 3:main.main [/source/go/scratch/scratch.go:28] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]

stacktraces.LongStackTrace(0):
0:runtime.Callers
/usr/local/go/src/runtime/extern.go:308
0x8e0a3
---
1:parasaurolophus/go/stacktraces.formatStackTrace
/source/go/stacktraces/stacktraces.go:235
0x8e050
---
2:parasaurolophus/go/stacktraces.LongStackTrace
/source/go/stacktraces/stacktraces.go:133
0x8daf7
---
3:main.main
/source/go/scratch/scratch.go:29
0x8e693
---
4:runtime.main
/usr/local/go/src/runtime/proc.go:267
0x431ab
---
5:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1197
0x6dc73
```
