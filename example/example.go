// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/stacktraces"
)

func main() {

	_, functionName, fileName, _, _ := stacktraces.FunctionInfo(nil)

	counter := 0

	loggerOptions := logging.LoggerOptions{
		BaseTags: []string{fileName},
		// Set base attributes using pointers when their values might change
		// over time.
		//                               |
		//                               V
		BaseAttributes: []any{"counter", &counter},
	}

	logger := logging.New(os.Stdout, &loggerOptions)

	defer func() {
		if r := recover(); r != nil {
			logger.Always(
				func() string {
					return fmt.Sprintf("recovered: %v", r)
				},
				logging.RECOVERED, r,
				logging.TAGS, []string{"PANIC"},
			)
		}
		counter += 1
		logger.Trace(
			func() string {
				return fmt.Sprintf("exiting %s", functionName)
			},
		)
	}()

	logger.Trace(func() string { return "you won't see this" })
	logging.SetVerbosity(logger, logging.TRACE)
	panic("deliberate")
}
