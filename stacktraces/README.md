_Copyright &copy; Kirk Rader 2024_

# Stack Traces for Logging and Error Messages in Go

Go's design is a throw-back (pun intended) to early 1970's era language design.
Its semantics can be summarized as "K&R C plus lexical closures and stack
unwinding protection, but lacking the expressive power of pointer arithmentic
due to being burdened with a garbage collector." One area particularly affected
by these design choices is error handling and reporting. Go has no mechanism for
throwing or catching exceptions. A further consequence is that Go programmers
are on their own for finding ways to include meaningful stack traces in logs as
an aid in debugging. (It would be nice if Go's `error` type implemented some
kind of automatic stack trace collection mechanism. That would require the
ability for custom `error` types to inherit such a mechanism and inheritance --
along with method overloading -- is another of the many incredibly useful
language features introduced into programming languages during the 1970's with
which Go has dispensed for, apparently, no better reason than that Go's
designers do not understand why they were introduced in the first place or else
believe that if a feature has ever been misused it must be excluded.)

This wrapper provides an admittedly somewhat painful work-around. It defines a
`stacktraces.StackTrace` struct which implements the `error` interface. Its
"constructor," `stacktraces.New(msg string, skipFrames any)`, captures a stack
trace at the time it is called using the same logic as the helper methods
`stacktraces.LongStackTrace(skipFrames any)` and
`stacktraces.ShortStackTrace(skipFrames any)`. 

The "skip frames" feature is supplied so that programmers can exclude the first
frames in call chains that always start with the stack frame formatting and
overall logging library implementation. There are three options supported for
the `skipFrames` parameter to the stack trace related functions:

- If passed a string, the stack trace will exclude all frames up to a frame for
  the function with the given name.

- If passed a non-negative number, the stack trace will omit the specified
  number of frames from the top of the call stack.

- If any other value, the stack trace will exclude all frames up to and
  including the function which created the trace.

For example, when called directly from `main.main`:

```go
// prints "main.main"
fmt.Println(stacktraces.FunctionName())

// prints a stack trace starting at "main.main"
fmt.Println(stacktraces.ShortStackTrace(-1))

// also prints a stack trace starting at "main.main"
fmt.Println(stacktraces.ShortStackTrace(stacktraces.FunctionName()))

// prints a stack trace starting at "runtime.main" (the caller of "main.main")
fmt.Println(stacktraces.ShortStackTrace("runtime.main"))
```

will print the following to `stdout`:

```
main.main
3:main.main [/source/go/scratch/scratch.go:29] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]
3:main.main [/source/go/scratch/scratch.go:32] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]
4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]
```

The second of these starts with `main.main` because it is the direct caller of
`stacktraces.ShortStackTrace(-1)` and the latter will exclude everything up to
and including its own frame when passed a negative `int`. Both
`stacktraces.LongStackTrace(any)` and `stacktraces.New(string, any)` behave
similarly. If a stack frame is not found matching the name passed to any of
these functions then the returned string will be empty. Similarly, if a positive
number is passed that is greater than the call stack depth the result will be
the empty string.

By far the most common usage patterns are to pass -1 when calling
`stacktraces.New(string, any)` directly, e.g. to create an `error` return value
that includes a stack trace or to pass the name of the calling function when
specifying a `stacktrace` attribute for a log entry. The
`stacktraces.FunctionName()` helper is provided to make this easy to implement in a
maintainable fashion. See [../example/example.go](../example/example.go) and
[stacktrace_test.go](./stacktrace_test.go) for examples of both these patterns.

The ability to specify arbitrary function names and skip-frame counts is
provided for completeness, i.e. to support largely hypothetical debugging
scenarios. In the real world you are likely only ever to find uses for passing
-1 or the result of calling `stacktraces.FunctionName()` and, possibly, 0 on
rare occasions, especially when you want to debug this library itself. If you
often find yourself using values for "skip frames" other than -1, 0 or the
result of `stacktraces.FunctionName()` you might want to consider why you think
that to be necessary and consider whether or not the resulting code might be
somewhat fragile. The implementation of `stacktraces.FunctionName()` is one
exception that demonstrates the problem. It is implemented by calling
`stacktraces.functionNameAt(int)`, passing a hard-coded constant which is the
number of stack frames it takes to reach the frame for
`stacktraces.FunctionName()`'s own caller. If the implementation of
`stacktraces.FunctionName()` ever changes such that there are more or fewer
intermediate functions on the stack, that hard-coded number would have to
change.

## Go Docs

```bash
$ go doc -all
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
```
