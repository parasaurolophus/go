// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

var (

	// An additional value that will be included in every log entry using
	// loggerOptions.BaseAttributes.
	counter = 0

	// Make main.main()'s source info available globally for logging.
	sourceInfo stacktraces.SourceInfo
)

// Writes two (not three) log entries to stdout and exits with status code 2 due
// to a deliberately caused panic.
//
//	{"time":"2024-03-15T04:31:00.30668337-05:00","verbosity":"TRACE","msg":"you will see this","counter":0,"file":{"function":"main.main","file":"/source/go/example/example.go","line":58},"tags":["EXAMPLE"]}
//	{"time":"2024-03-15T04:31:00.306936205-05:00","verbosity":"ALWAYS","msg":"recovered: runtime error: invalid memory address or nil pointer dereference","counter":1,"recovered":"runtime error: invalid memory address or nil pointer dereference","file":{"function":"runtime.panicmem","file":"/usr/local/go/src/runtime/panic.go","line":261},"stacktrace":"5:main.finally [/source/go/example/example.go:96] < 6:runtime.gopanic [/usr/local/go/src/runtime/panic.go:770] < 7:runtime.panicmem [/usr/local/go/src/runtime/panic.go:261] < 8:runtime.sigpanic [/usr/local/go/src/runtime/signal_unix.go:881] < 9:main.uncheckedPointer [/source/go/example/example.go:70] < 10:main.main [/source/go/example/example.go:65] < 11:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 12:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC"]}
func main() {

	_, sourceInfo, _ = stacktraces.FunctionInfo(nil)

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

	// Create an asynchronous logger.
	logger := logging.New(os.Stdout, &loggerOptions)

	// Log every exit from main.
	defer finally(logger)

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
	uncheckedPointer(nil)
}

// Called from main() to demonstrate logging panics.
func uncheckedPointer(p *int) {
	fmt.Println(*p)
}

// Invoked by main() using defer to ensure logging on exit, even if a panic
// occurs.
func finally(logger *logging.Logger) {

	// Check to see if a panic occurred.
	r := recover()

	if r == nil {

		// Exiting normally, so log at TRACE level.
		logger.Trace(
			func() string {
				return fmt.Sprintf("exiting %s", sourceInfo.Function)
			})

	} else {

		// Always log a stack trace, source info, and recovered value when a
		// panic occurs.

		// Increment counter.
		counter += 1

		// Note that counter is 1 in this entry.
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

	// Signal an unsuccessful exit status when a panic has occurred.
	if r != nil {
		os.Exit(2)
	}
}
