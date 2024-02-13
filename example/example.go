// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

var (

	// Specify the logger configuration.
	loggerOptions = logging.LoggerOptions{
		BaseTags: []string{"EXAMPLE"},
	}

	// Create a logger.
	logger = logging.New(os.Stdout, &loggerOptions)

	// Conventional tags when logging a panic.
	panicTags = []string{"PANIC", "ERROR", "SEVERE"}
)

func main() {

	// enable most verbose logging
	logger.SetVerbosity(logging.TRACE)

	// for use as the parameter to a log entry's "stacktrace" attribute
	functionName := stacktraces.FunctionName()

	defer logger.Defer(

		// panics in main.main should cause abnormal termination
		true,

		// no clean-up required for main.main since sender() will close the
		// channel passed to it
		nil,

		// remaining parameters would be passed to logger.Always() if recover()
		// returned non-nil in main.main

		func(r any) string {
			return fmt.Sprintf("%s recovered from '%v'", functionName, r)
		},

		// add conventional attributes for logging a panic
		logging.STACKTRACE, functionName,
		logging.TAGS, panicTags,
	)

	// make a channel for receiving values from a goroutine
	ch := make(chan int)

	logger.Trace(func() string { return fmt.Sprintf("%s starting goroutine", functionName) })
	go sender(ch)

	logger.Trace(func() string { return fmt.Sprintf("%s consuming output from goroutine", functionName) })

	// consume values from the goroutine until ch is closed
	for v := range ch {

		fmt.Println(v)
	}

	logger.Trace(func() string { return fmt.Sprintf("%s exiting normally", functionName) })
}

// Demonstrate logging from a goroutine, including logging and recovering from a
// panic.
func sender(ch chan int) {

	functionName := stacktraces.FunctionName()

	defer logger.Defer(

		// log but otherwise ignore panics in this goroutine
		false,

		// clean-up function is always called when this function exits
		func() {
			logger.Trace(func() string { return fmt.Sprintf("%s closing channel", functionName) })
			close(ch)
		},

		// remaining arguments are passed to logger.AlwaysContext() when
		// recover() returns non-nil

		func(recovered any) string {
			// return nil as second value when panicing in a goroutine so that
			// others can complete normally
			return fmt.Sprintf("%s recovered from '%v'", functionName, recovered)
		},

		// include conventional attributes for logging panics
		logging.STACKTRACE, functionName,
		logging.TAGS, panicTags,
	)

	// send values to ch's consumer, triggering a panic at some point along the
	// way
	for v := 0; v < 10; v++ {

		if v > 4 {
			// deliberately trigger a panic to demonstrate logger.Defer()
			logger.Trace(func() string { return "sender deliberately causing a panic" })
			panic("deliberate panic")
		}

		ch <- v
	}
}
