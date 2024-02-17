_Copyright &copy; Kirk Rader 2024_

# Stack Traces for Logging and Error Messages in Go

The primary documentation for this package is inline comments in the source
code, viewable by way of `go doc`.

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

## Examples

```go
package main

import (
	"fmt"
	"parasaurolophus/go/stacktraces"
)

func main() {
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

	// as are stacktraces.FunctionName() and stacktraces.FunctionInfo() helpers
	fmt.Println(stacktraces.FunctionName())
	frame, functionName, fileName, lineNumber, found := stacktraces.FunctionInfo(-1)
	fmt.Printf(
		"Info about stacktraces.FunctionInfo()'s caller's caller:\n"+
			"frame number: %d\n"+
			"function name: '%s'\n"+
			"file name: '%s'\n"+
			"line number: %d\n"+
			"found: %v\n",
		frame, functionName, fileName, lineNumber, found)
}
```

writes the following to `stdout`:

```
err.Error(): message

stackTrace.ShortTrace(): 3:main.main.func1 [/source/go/scratch/scratch.go:14] < 4:main.main [/source/go/scratch/scratch.go:15] < 5:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 6:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]

stackTrace.LongTrace():
3:main.main.func1
/source/go/scratch/scratch.go:14
0x947cb
---
4:main.main
/source/go/scratch/scratch.go:15
0x947b0
---
5:runtime.main
/usr/local/go/src/runtime/proc.go:271
0x4560b
---
6:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1222
0x739d3

stacktraces.ShortStackTrace(0): 0:runtime.Callers [/usr/local/go/src/runtime/extern.go:325] < 1:parasaurolophus/go/stacktraces.formatStackTrace [/source/go/stacktraces/stacktrace.go:230] < 2:parasaurolophus/go/stacktraces.ShortStackTrace [/source/go/stacktraces/stacktrace.go:68] < 3:main.main [/source/go/scratch/scratch.go:27] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]

stacktraces.LongStackTrace(0):
0:runtime.Callers
/usr/local/go/src/runtime/extern.go:325
0x94413
---
1:parasaurolophus/go/stacktraces.formatStackTrace
/source/go/stacktraces/stacktrace.go:230
0x943c0
---
2:parasaurolophus/go/stacktraces.LongStackTrace
/source/go/stacktraces/stacktrace.go:59
0x93ad7
---
3:main.main
/source/go/scratch/scratch.go:28
0x949c7
---
4:runtime.main
/usr/local/go/src/runtime/proc.go:271
0x4560b
---
5:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1222
0x739d3
main.main
Info about stacktraces.FunctionInfo()'s caller's caller:
frame number: 4
function name: 'runtime.main'
file name: '/usr/local/go/src/runtime/proc.go'
line number: 271
found: true
```

### skipFrames

Most of the functions exposed by this library take a parameter named
`skipFrames` of type `any`. Such a parameter is used to specify which frames in
the current call stack are considered matches, i.e. included in a stack trace or
selected as the frame returned by `stacktraces.New()`,
`stacktraces.LongStackTrace()`, `stacktraces.ShortStackTrace()`,
`stacktraces.FunctionName()` and `stacktraces.FunctionInfo()`. In all cases,
`skipFrames` is interpreted as follows:

| `skipFrames`             | Selected Frame                                                                                                  |
|--------------------------|-----------------------------------------------------------------------------------------------------------------|
| `int` >= 0               | Skip the given number of frames from the beginning of the current call stack                                    |
| `int` < 0                | Skip `n` frames after the frame for the currently executing function where `n` is the magnitude of `skipFrames` |
| `string`                 | Skip all frames preceding the function with the given name                                                      |
| `nil` or any other value | Skip all frames preceding the currently executing function                                                      |

To further elucidate:

- If `skipFrames` is an `int` and `skipFrames >= 0`, the entire call stack is
  matched after skipping the specified number of frames.

- If `skipFrames` is an `int` and `skipFrames < 0`, the entire call stack is
  matched starting with the specified number of frames after the function which
  called the `stacktraces` function, as described below in more detail.

- If `skipFrames` is a `string`, the entire call stack is matched starting with
  the function with the specified name.

- If `skipFrames` is any other value, including `nil`, the entire call stack is
  matched starting with the function which called the `stacktraces` function,
  again as described below in more detail.

For example,

```go
frameNumber, functionName, _, _, found := stacktraces.FunctionInfo(0)
fmt.Printf(
    "%d, '%s', %v\n",
    frameNumber,
    functionName,
    found)
```

always prints

```
0, 'runtime.Callers', true
```

because, due to the way that `stacktraces.FunctionInfo()` is implemented,
`runtime.Callers` is always the first frame on the call stack.

Similarly,

```go
package main

import (
	"fmt"
	"parasaurolophus/go/stacktraces"
)

func main() {
	frameNumber, functionName, _, _, found := stacktraces.FunctionInfo(nil)
	fmt.Printf(
		"%d, '%s', %v\n",
		frameNumber,
		functionName,
		found)
}
```

prints

```
3, 'main.main', true
```

I.e. passing `nil` for `skipFrames` (or any value other than a `string` or
`int`) will select the frame for the currently executing function, from that
function's point of view.

In the preceding example, note that the frame number for "the currenly executing
function" is 3 because of how this library is structured. Thus, passing 3 or
`nil` will for `skipFrames` will produce the same result. I.e.

```go
frameNumber, functionName, _, _, found := stacktraces.FunctionInfo(3)
fmt.Printf(
    "%d, '%s', %v\n",
    frameNumber,
    functionName,
    found)
```

prints

```
3, 'main.main', true
```

Other non-negative `int` values simply skip the specified number of frames. So
passing 4 will select the currently executing function's caller, passing 5 will
select that function's caller, and so on. If you pass a number that is greater
than the current calling stack depth then the last return value from
`stacktraces.FunctionInfo()` will be `false` and the others will be set to their
"zero values" (in Go's terminology). I.e.

```go
frameNumber, functionName, _, _, found := stacktraces.FunctionInfo(100)
fmt.Printf(
    "%d, '%s', %v\n",
    frameNumber,
    functionName,
    found)
```

prints

```
0, '', false
```

Passing a `string` for `skipFrames` searches for a frame with the specified
function name. I.e.

```go
package main

import (
	"fmt"
	"parasaurolophus/go/stacktraces"
)

func main() {
	func() {
		name := stacktraces.FunctionName()
		frameNumber, functionName, _, _, found := stacktraces.FunctionInfo(name)
		fmt.Printf(
			"%d, '%s', %v\n",
			frameNumber,
			functionName,
			found)
	}()
}
```

prints

```
3, 'main.main.func1', true
```

because `stacktraces.FunctionName()` returns `main.main.func1` for the anonymous
function in `main.main` and `stacktraces.FunctionInfo("main.main.func1")`
selects the frame for the function with that name. This behaves similarly to the
case of passing too large an `int` if no function is found on the current call
stack with the specified name. I.e.

```go
frameNumber, functionName, _, _, found := stacktraces.FunctionInfo("noSuchFunction")
fmt.Printf(
    "%d, '%s', %v\n",
    frameNumber,
    functionName,
    found)
```

prints

```
0, '', false
```

Passing a negative `int` behaves like a combination of passing `nil` and a
positive `int`. Specifically, it subtracts (i.e. adds the magnitude of) the
numeric value of `skipFrames` to the same offset used when `skipFrames` is
neither an `int` nor `string`. This has the effect of selecting the frame that
is `n` frames beyond the frame for the currently executing function, where `n =
|skipFrames|`. I.e.

```go
package main

import (
	"fmt"
	"parasaurolophus/go/stacktraces"
)

func main() {
	func() {
		frameNumber, functionName, _, _, found := stacktraces.FunctionInfo(-1)
		fmt.Printf(
			"%d, '%s', %v\n",
			frameNumber,
			functionName,
			found)
	}()
}
```

prints

```
4, 'main.main', true
```

because `main.main` is the caller of `main.main.func1`.

The same rules apply when passing a `skipFrames` parameter to
`stacktraces.New()`, `stacktraces.LongStackTrace()` and
`stacktraces.ShortStackTrace()`. In these cases the trace will begin with the
selected frame and include all following frames. E.g.

```go
package main

import (
	"fmt"
	"parasaurolophus/go/stacktraces"
)

func main() {
	fmt.Print(stacktraces.LongStackTrace(nil))
}
```

prints

```
3:main.main
/source/go/scratch/scratch.go:11
0x93d83
---
4:runtime.main
/usr/local/go/src/runtime/proc.go:271
0x4560b
---
5:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1222
0x739d3
```

The most common usage for most debugging and error reporting scenarios is to
pass `nil` as the `skipFrames` parameter to any of the preceding functions, but
the other values are supported for special purposes, such as embedding stack
traces in log entries (see [../logging/README.md](../logging/README.md)) or for
debugging this library, itself.
