// Copyright Kirk Rader 2024

package main

import (
	"context"
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

func main() {

	type (
		errorCounters struct {
			Error1 int `json:"error1"`
			Error2 int `json:"error2"`
		}
	)

	counters := errorCounters{}

	// display the name of the currenly executing function, main.main
	fmt.Printf("FunctionName(): %s\n\n", stacktraces.FunctionName())

	// display a one-line stack trace starting at the frame for runtime.main,
	// which calls main.main within the Go runtime
	fmt.Printf("ShortStackTrace(\"runtime.main\"): %s\n\n", stacktraces.ShortStackTrace("runtime.main"))

	// display a one-line stack trace starting at the currently executing
	// function
	fmt.Printf("ShortStackTrace(-1): %s\n\n", stacktraces.ShortStackTrace(-1))

	// display a one-line stack trace starting at the top frame, which is always
	// runtime.Callers
	fmt.Printf("ShortStackTrace(0): %s\n\n", stacktraces.ShortStackTrace(0))

	// use a stacktraces.StackTrace as an error instance
	err := func() error {
		return stacktraces.New("StackTrace as error", -1)
	}()

	if err != nil {
		stackTrace := err.(stacktraces.StackTrace)
		fmt.Printf("stackTrace.Error(): %s\nstackTrace.LongTrace():\n%s\n", stackTrace.Error(), stackTrace.LongTrace())
	}

	// use default options for everything except BaseAttributes and BaseTags
	options := logging.LoggerOptions{
		// use a reference to counters in base attributes
		BaseAttributes: []any{"counters", &counters},
		BaseTags:       []string{"example"},
	}

	// construct a logging.Logger
	logger := logging.New(os.Stdout, &options)

	ctx := context.Background()

	// arrange to log a stack trace when panicing
	defer logger.OnPanic(
		ctx,
		func(r any) (string, any) {
			// by returning r as second argument, OnPanic will re-invoke panic
			return fmt.Sprintf("panic: %#v", r), r
		},
		logging.TAGS, []string{"ERROR", "PANIC", "SEVERE"},
		logging.STACKTRACE, nil,
	)

	// TRACE is disabled by default
	logger.Trace(
		ctx,
		func() string {
			// lazy evaluation suppresses the invocation of message builder
			fmt.Println("you won't see this")
			return "fail"
		})

	// enable TRACE and log again
	logger.SetVerbosity(logging.TRACE)

	n := 42

	logger.Trace(
		ctx,
		func() string {
			fmt.Println("you will see this")
			// note that message builder closures have access to the lexical
			// environment in which they are created
			return fmt.Sprintf("n is %d", n)
		},
		logging.STACKTRACE, nil)

	fmt.Println()
	counters.Error1 += 1

	logger.Always(
		ctx,
		func() string {
			return "note that counters in the log entry reflects the current value of error1"
		},
		logging.TAGS, []string{"reference_test"})

	// deliberately panic
	testPanic := func() {

		fmt.Printf("\ndeliberately panicing\n")
		panic("example")
	}

	testPanic()

	// execution won't reach here due to the panic
	logger.Always(
		ctx,
		func() string {
			return "you won't see this"
		},
		logging.TAGS, []string{"panic_test"})

}
