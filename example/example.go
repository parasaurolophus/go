// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
	"strconv"
)

var (
	// The number of values sent by a goroutine.
	sent = 0
	// The number of values received from the goroutine.
	received = 0

	loggerOptions = logging.LoggerOptions{
		BaseTags: []string{"EXAMPLE"},
		BaseAttributes: []any{
			"sent", &sent,
			"received", &received,
		},
	}

	logger = logging.New(os.Stdout, &loggerOptions)
)

// Print the number of values sent by and received from a goroutine to stdout.
//
// Logging verbosity defaults to OPTIONAL but may be set using a command-line
// argument.
func main() {
	verbosity := logging.OPTIONAL
	logger.SetVerbosity(verbosity)
	panicInMain := false
	panicAgain := true
	parseArg(1, "logging.Verbosity", &verbosity)
	parseArg(2, "bool", &panicInMain)
	parseArg(3, "bool", &panicAgain)
	fmt.Printf("\nverbosity: %s, panicInMain: %v, panicAgain: %v\n\n", verbosity, panicInMain, panicAgain)
	logger.SetVerbosity(verbosity)
	functionName := stacktraces.FunctionName()
	defer logger.Finally(
		panicAgain,
		func() {
			logger.Optional(nil, logging.TAGS, "DEBUG")
		},
		func(r any) string {
			return fmt.Sprintf("%s panicing: %v", functionName, r)
		},
		logging.STACKTRACE, functionName,
		logging.TAGS, []string{"PANIC", "SEVERE"})

	ch := make(chan int)
	logger.Trace(
		func() string { return fmt.Sprintf("%s starting sender goroutine", functionName) },
		logging.TAGS, "DEBUG")
	go sender(ch)
	logger.Trace(
		func() string { return fmt.Sprintf("%s consuming output from sender goroutine", functionName) },
		logging.TAGS, "DEBUG")
	for v := range ch {
		received += 1
		logger.Fine(func() string { return strconv.Itoa(v) })
	}
	fmt.Printf("\n%d sent, %d received\n\n", sent, received)
	if panicInMain {
		panic("another deliberate panic")
	}
	logger.Trace(func() string { return fmt.Sprintf("%s exiting normally", functionName) })
}

// Goroutine that sends int values to a channel.
//
// This deliberately panics after sending a few values as a demonstration of
// logging.Logger.Defer().
func sender(ch chan int) {
	functionName := stacktraces.FunctionName()
	defer logger.Finally(
		false,
		func() {
			logger.Trace(func() string { return fmt.Sprintf("%s goroutine closing channel", functionName) })
			close(ch)
		},
		func(recovered any) string {
			return fmt.Sprintf("%s recovered from: %v", functionName, recovered)
		},
		logging.STACKTRACE, functionName,
		logging.TAGS, []string{"PANIC", "MEDIUM"},
	)
	for v := 0; v < 10; v++ {
		if v > 4 {
			logger.Trace(func() string { return "sender deliberately causing a panic" })
			panic("deliberate panic")
		}
		sent += 1
		ch <- v
	}
}

func parseArg(index int, typeName string, val any) {
	if len(os.Args) <= index {
		logger.Fine(
			func() string { return fmt.Sprintf("optional argument %d (of type %s) not supplied", index, typeName) },
			logging.TAGS, "DEBUG")
		return
	}
	n, err := fmt.Sscan(os.Args[index], val)
	if err != nil {
		logger.Optional(
			func() string { return err.Error() },
			logging.STACKTRACE, nil,
			logging.TAGS, []string{"ERROR", "USER", "BAD_ARGS"},
		)
		os.Exit(1)
	}
	if n != 1 {
		logger.Optional(
			func() string { return fmt.Sprintf("expected 1 %s, got %d", typeName, n) },
			logging.STACKTRACE, nil,
			logging.TAGS, []string{"ERROR", "USER", "BAD_ARGS"},
		)
		os.Exit(1)
	}
	logger.Trace(
		func() string { return fmt.Sprintf("parsed arg %d: %v", index, val) },
		logging.TAGS, "DEBUG")
}
