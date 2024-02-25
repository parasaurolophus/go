_Copyright &copy; Kirk Rader 2024_

# Concurrency in Go

## Channels and goroutines

Go has a built-in, lightweight concurrency mechanism called "goroutines." (Like
"co-routines." Don't forget to tip the wait staff!) Any function invocation can
be turned into a concurrently executing goroutine using the `go` statement.

```go
package main

import (
	"fmt"
)

// Prints the following to stdout (in a particularly round-about way):
//
//	1
//	2
//	3
func main() {

	// make a channel for communicating between goroutines
	ch := make(chan int)

	// launch a goroutine that will emit each of its parameters to ch before
	// closing it
	go func(parameters ...int) {

		// ensure that the ch is closed whether or not this goroutine terminates
		// normally
		defer func() { close(ch) }()

		// emit each parameter to ch
		for _, n := range parameters {
			ch <- n
		}

	}(1, 2, 3)

	// consume all value produced by ch until it is closed by the goroutine
	for m := range ch {

		fmt.Println(m)
	}
}
```

Both `main.main()` and the function invoked using the `go` keyword continue to
execute in parallel. In this case, they use a `channel` to communicate and
synchronise their execution. The Go documentation makes much of this pattern,
even basing the name of the language on the `go` statement's keyword and
adopting as a motto "don't communicate by sharing memory, share memory by
communicating." (It is worth noting that as worthwhile as that sentiment may be,
it is utterly impossible to achieve in Go or any other language. In particular,
goroutines that communicate using channels share memory in the form of the
channel object, itself. But I digress....) Go also supports more traditional
thread synchronization mechanisms such as mutexes for use with goroutines, but
using channels in a pattern similar to that shown in the preceding example often
results in cleaner code that is less prone to deadlocks and race conditions.

## Handling Errors in goroutines

The first example, shown above, has a number of features worth examining in a
bit more detail. Note that the `ch` is created in the main goroutine and closed
in the worker goroutine. This is part of the synchronization pattern that ensure
the goroutine will be allowed to complete before the program exits. Had
`main.main()` not waited for `ch` to close using the `for v := range ch { ... }`
idiom, the goroutine would have been killed if it hadn't completed its work
before `main.main()` returned.

But that means the example program would be in danger of hanging if the
goroutine ever terminated without closing `ch`. The anonymous function protects
against that possibility by calling `close(ch)` in a `defer` handler. Like the
`go` keyword, `defer` statements arrange for functions to be invoked in a
special context. Deferred functions are always called when execution leaves the
execution context in which they are declared even when errors occur, like
`ensure` statements in Ruby or `try ... finally ...` in C++, Java, JavaScript
etc.

In particular, if the body of the goroutine in this example were ever to cause a
`panic`, the deferred invocation of `close(ch)` would still occur. That would be
of little comfort to `main.main()`, unfortunately, because the panic would kill
the whole program despite its having at least closed the channel. A bit like
burning down the barn after locking its door once the horses have escaped.

Since goroutines are often used in contexts where inversion of control or
similar patterns are in use, it is often desirable to catch errors that would
otherwise cause `main.main()` to terminate. The combination of built-in `defer`,
`panic()` and `recover()` allow for exactly such scenarios. A slight
modification to the preceding example demonstrates the problen (note the
invocation of `panic()` after the first value is emitted on `ch`):

```go
package main

import (
	"fmt"
)

func main() {

	ch := make(chan int)

	go func(parameters ...int) {

		defer func() { close(ch) }()

		for _, n := range parameters {
			ch <- n
			panic("deliberate")
		}

	}(1, 2, 3)

	for m := range ch {

		fmt.Println(m)
	}
}
```

which produces the following output when run:

```
panic: deliberate

goroutine 18 [running]:
main.main.func1({0x400003c7b0?, 0x0?, 0x0?})
	/source/go/scratch/scratch.go:30 +0x78
created by main.main in goroutine 1
	/source/go/scratch/scratch.go:21 +0xac
exit status 2
```

But see what happens when the goroutine gracefully handles the panic:

```go
package main

import (
	"fmt"
	"os"
)

func main() {

	ch := make(chan int)

	go func(parameters ...int) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "recovered: %v\n", r)
			}
			close(ch)
		}()

		for _, n := range parameters {
			ch <- n
			panic("deliberate")
		}

	}(1, 2, 3)

	for m := range ch {

		fmt.Println(m)
	}
}
```

produces the following output and exists cleanly, thanks to the use of
`recover()` in the deferred clean-up function.

```
recovered: deliberate
1
```

Note that the logging library provides `Logger.Finally()` and
`Logger.FinallyContext()` methods that are intended for use with `defer` to
ensure that panics are properly logged and optionally recovered.
