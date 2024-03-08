// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

// Called from main() to demonstrate logging panics.
func deliberatePanic() {
	panic("deliberate")
}

// Writes three (not four) log entries to stdout and exits with status code 1
// due to a deliberately caused panic.
//
// {"time":"2024-03-08T08:18:15.535425196-06:00","verbosity":"TRACE","msg":"you will see this","counter":0,"file":{"function":"main.main","file":"/source/go/example/example.go","line":116},"tags":["EXAMPLE"]}
// {"time":"2024-03-08T08:18:15.536047625-06:00","verbosity":"ALWAYS","msg":"recovered: deliberate","counter":1,"recovered":"deliberate","file":{"function":"main.deliberatePanic","file":"/source/go/example/example.go","line":14},"stacktrace":"5:main.main.func1 [/source/go/example/example.go:81] < 6:runtime.gopanic [/usr/local/go/src/runtime/panic.go:770] < 7:main.deliberatePanic [/source/go/example/example.go:14] < 8:main.main [/source/go/example/example.go:123] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC"]}
// {"time":"2024-03-08T08:18:15.536090403-06:00","verbosity":"TRACE","msg":"exiting main.main","counter":1,"tags":["EXAMPLE"]}
func main() {

	// Get the calling function's name -- main.main in this case.
	functionName := stacktraces.FunctionName()

	// An additional value that will be included in every log entry using
	// loggerOptions.BaseAttributes.
	counter := 0

	// Configuration options for a logging.Logger.
	loggerOptions := logging.LoggerOptions{

		// Include EXAMPLE as a tag in every log entry.
		BaseTags: []string{"EXAMPLE"},

		// Set base attributes using pointers when their values might change
		// over time -- but then beware of race conditions when using
		// asynchronous loggers. Use a synchronized function where appropriate.
		//                                   |
		//                                   V
		BaseAttributes: []any{"counter", &counter},
	}

	// Create a synchronous logger.
	logger := logging.New(os.Stdout, &loggerOptions)

	//////////////////////////////////////////////////////////////////////////
	// ALWAYS log panics with stack trace & recovered attributes in addition to
	// TRACE logging every exit from main().
	defer func() {

		// Check to see if a panic occurred.
		r := recover()

		// Always log a stack trace, source info, and recovered value when a
		// panic occurs.
		if r != nil {

			// Increment counter.
			counter += 1

			// Note that counter is 1 in this entry and subsequent entry.
			logger.Always(
				func() string {
					return fmt.Sprintf("recovered: %v", r)
				},
				logging.STACKTRACE, nil,
				logging.RECOVERED, r,
				logging.TAGS, []string{logging.PANIC},
				logging.FILE, logging.FILE_SKIPFRAMES_FOR_PANIC,
			)
		}

		logger.Trace(
			func() string {
				return fmt.Sprintf("exiting %s", functionName)
			})

		// logger.Stop is a no-op for synchronous logger, but called here as
		// insurance if we ever decide to change to asyncrhonous logging.
		logger.Stop()

		// Signal an unsuccessful status when a panic occurs.
		if r != nil {
			os.Exit(2)
		}
	}()
	// end of defer
	//////////////////////////////////////////////////////////////////////////

	// Default verbosity is FINE, so logger.Trace() does not write an entry.
	logger.Trace(func() string { return "you won't see this" })

	// Adjust verbosity level
	logger.SetVerbosity(logging.TRACE)

	// Now logger.Trace() is included in output. Note that counter is 0 in this
	// entry.
	logger.Trace(
		func() string {
			return "you will see this"
		},
		logging.FILE, logging.FILE_SKIPFRAMES_FOR_CALLER)

	// Deliberately cause a panic so as to demonstrate deferred logging.
	deliberatePanic()
}
