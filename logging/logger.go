// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"io"
	"log/slog"
)

type (

	// Type of function passed to logging methods for lazy evaluation of message
	// formatting.
	//
	// The returned string becomes the value of the log entry's msg attribute.
	//
	// Such a function is invoked only if a given verbosity is enabled for a
	// given logger.
	MessageBuilder func() string

	Logger interface {
		Trace(context.Context, MessageBuilder, ...any)
		Fine(context.Context, MessageBuilder, ...any)
		Optional(context.Context, MessageBuilder, ...any)
		Always(context.Context, MessageBuilder, ...any)
		Enabled(context.Context, Verbosity) bool
		BaseAttributes() []any
		SetBaseAttributes(...any)
		BaseTags() []string
		SetBaseTags(...string)
		Verbosity() Verbosity
		SetVerbosity(Verbosity)
	}
)

// Specially handled attributes.
const (

	// Conventional attribute for including source file information in a log
	// entry.
	FILE = "file"

	// syncLogger.log() will include the value returned by recover() when
	// logging a panic.
	RECOVERED = "recovered"

	// Values of "stacktrace" attributes will be replaced with one-line stack
	// traces for the function that called the given logging method.
	STACKTRACE = "stacktrace"

	// Value will be merged with the currently configured
	// LoggerOptions.BaseTags.
	TAGS = "tags"

	// Value to use for FILE attribute when logging normal functionality.
	FILE_SKIPFRAMES_FOR_CALLER = -2

	// Value to use for FILE attribute when logging a panic.
	FILE_SKIPFRAMES_FOR_PANIC = -4
)

// Returns a newly created, wrapped instance of slog.Logger.
//
// Log entries written using the returned Logger instance will have "verbosity"
// attributes instead of "level" attributes and the values of their "stacktrace"
// attributes, if present, will be replaced as if by an invocation of
// ShortStackTrace(skipFrames) where skipFrames is the value of  the
// "stacktrace" attribute passed to a logging method. The final set of
// attributes for each log entry will be the result of combining the value of
// LoggerOptions.BaseAttributes and LoggerOptions.BaseTags with the attributes
// passed to the given logging method.
//
// For example:
//
//	type Counters struct {
//	    Error1 int `json:"error1"`
//	    Error2 int `json:"error2"`
//	}
//
//	counters := Counters{}
//
//	options := logging.LoggerOptions{
//	    BaseAttributes: []any{"counters", &counters},
//	    BaseTags:       []string{"foo", "bar"},
//	}
//
//	logger := logging.New(os.Stdout, &options)
//	n := 42
//	counters.Error1 += 1
//
//	logger.Optional(
//	    func() string { return fmt.Sprintf("n = %d", n) },
//	    logging.STACKTRACE, nil,
//	    logging.TAGS, []string{"hoo"},
//	    "baz", "waka")
//
// produces a log entry like:
//
//	{"time":"2024-02-11T06:16:41.852302853-06:00","verbosity":"OPTIONAL","msg":"n = 42","counters":{"error1":1,"error2":0},"baz":"waka","stacktrace":"5:main.main [/source/go/scratch/scratch.go:29] < 6:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 7:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","tags":["foo","bar","hoo"]}
//
// Note that if LoggerOptions.ReplaceAttr is not nil, it will be called as
// described by the documentation for slog.HandlerOptions.ReplaceAttr indirectly
// through a custom replacer function that replaces "level" with "verbosity" as
// just described.
//
// Note also that the values for attributes in LoggerOptions.BaseAttributes may
// be passed by value or reference. Passing by reference allows for cases where
// each log entry should include the current value for that attribute rather
// than a copy of the value at the time the Logger was created.
func New(writer io.Writer, options *LoggerOptions) Logger {

	l := new(loggerStruct)
	if options != nil {
		l.options = *options
	}

	if l.options.Level == nil {
		l.options.Level = new(slog.LevelVar)
	}

	handlerOptions := new(slog.HandlerOptions)
	handlerOptions.AddSource = false
	handlerOptions.Level = l.options.Level
	handlerOptions.ReplaceAttr = newAttrReplacer(l.options.ReplaceAttr)
	l.wrapped = slog.New(slog.NewJSONHandler(writer, handlerOptions))
	return l
}
