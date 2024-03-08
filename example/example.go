// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

// Writes three (not four) log entries to stdout and exits with status code 1
// due to a deliberately caused panic.
//
//	{"time":"2024-03-08T04:00:47.671436247-06:00","verbosity":"TRACE","msg":"you will see this","counter":0,"tags":["EXAMPLE"]}
//	{"time":"2024-03-08T04:00:47.67194538-06:00","verbosity":"ALWAYS","msg":"recovered: deliberate","counter":1,"recovered":"deliberate","stacktrace":"5:main.main.func1 [/source/go/example/example.go:54] < 6:runtime.gopanic [/usr/local/go/src/runtime/panic.go:770] < 7:main.main [/source/go/example/example.go:91] < 8:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC"]}
//	{"time":"2024-03-08T04:00:47.671993582-06:00","verbosity":"TRACE","msg":"exiting main.main","counter":1,"tags":["EXAMPLE"]}
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
		// asynchronous loggers.         |
		//                               V
		BaseAttributes: []any{"counter", &counter},
	}

	logger := logging.New(os.Stdout, &loggerOptions)

	defer func() {

		// Check to see if a panic occurred.
		r := recover()

		// Always log a stack trace and the recovered value when a panic occurs.
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
			)
		}

		logger.Trace(
			func() string {
				return fmt.Sprintf("exiting %s", functionName)
			},
		)

		// logger.Stop is a no-op for synchronous logger, but called here as
		// insurance if we ever decide to change to asyncrhonous logging.
		logger.Stop()

		// Signal an unsuccessful status when a panic occurs.
		if r != nil {
			os.Exit(2)
		}
	}()

	// Default verbosity is FINE, so logger.Trace() does not write an entry.
	logger.Trace(func() string { return "you won't see this" })

	// Adjust verbosity level
	logger.SetVerbosity(logging.TRACE)

	// Now logger.Trace() is included in output. Note that counter is 0 in this
	// entry.
	logger.Trace(func() string { return "you will see this" })

	// Deliberately cause a panic so as to demonstrate deferred logging.
	panic("deliberate")
}
