// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"io"
	"log/slog"
)

type (

	// Verbosity-based nomenclature used in place of slog.Level.
	Verbosity int

	// Type of function passed to logging methods for lazy evaluation of message
	// formatting.
	//
	// The returned string becomes the value of the log entry's msg attribute.
	//
	// Such a function is invoked only if a given verbosity is enabled for a
	// given logger.
	MessageBuilder func() string

	// Type of function passed to Logger.OnPanic().
	//
	//  - r receives the value returned from recover()
	//  - first return value is the message to log
	//  - second return value should be r or nil
	//
	// Logger.OnPanic() will re-invoke panic() if the second return value is not
	// nil.
	PanicHandler func(r any) (string, any)

	// Configuration parameters for an instance of Logger.
	LoggerOptions struct {

		// Pass through to HandlerOptions for the wrapped slog.Logger.
		AddSource bool

		// Initial set of attributes that will be added to every log entry.
		BaseAttributes []any

		// Initial set of tags that will be added to every log entry.
		BaseTags []string

		// Shared slog.LevelVar, if desired; a Leveler will be created if this
		// is nil.
		Level *slog.LevelVar

		// If not nil, an attribute replacer function that will be called in
		// addition to replacing "level" attributes with "verbosiy" and other
		// special attribute handling.
		ReplaceAttr func([]string, slog.Attr) slog.Attr
	}

	// Wrapper for an instance of slog.Logger.
	Logger struct {

		// This logger's configuration.
		options LoggerOptions

		// The wrapped slog.Logger.
		wrapped *slog.Logger
	}
)

// Mapping of Verbosity to slog.Level values.
const (

	// Only emit a log entry when the most verbose output is specified.
	TRACE = Verbosity(slog.LevelDebug)

	// Only emit a log entry when unusually verbose output is specified.
	FINE = Verbosity(slog.LevelInfo)

	// Only emit a log entry when conventionally verbose output is specified.
	OPTIONAL = Verbosity(slog.LevelWarn)

	// Always emit a log entry.
	ALWAYS = Verbosity(slog.LevelError)
)

// Specially handled attributes
const (

	// Value of "stacktrace" attributes will be replaced with one-line stack
	// traces for the function that called the given logging method.
	//
	// See stacktraces.ShortStackTrace(any)
	STACKTRACE = "stacktrace"

	// Value will be merged with the currently configured base tags.
	TAGS = "tags"
)

// Returns a newly created, wrapped instance of slog.Logger.
//
// Log entries written using the returned Logger instance will have "verbosity"
// attributes instead of "level" attributes and the values of their "stacktrace"
// attributes, if present, will be replaced as if by an invocation of
// ShortStackTrace(caller) where caller is the name of the function that calls a
// logging method. The final set of attributes for each log entry will be the
// result of combining the value of LoggerOptions.BaseAttributes and
// LoggerOptions.BaseTags with the attributes passed to the given logging
// method.
//
// For example:
//
//	type Counters struct {
//		Error1 int `json:"error1"`
//		Error2 int `json:"error2"`
//	}
//
//	counters := Counters{}
//
//	options := logging.LoggerOptions{
//		BaseAttributes: []any{"counters", &counters},
//		BaseTags:       []string{"foo", "bar"},
//	}
//
//	logger := logging.New(os.Stdout, &options)
//
//	n := 42
//
//	logger.Optional(
//		context.Background(),
//		func() string { return fmt.Sprintf("n = %d", n) },
//		logging.STACKTRACE, nil,
//		logging.TAGS, []string{"hoo"},
//		"baz", "waka")
//
// produces a log entry like:
//
//	{"time":"2024-02-09T06:56:10.285166661-06:00","verbosity":"OPTIONAL","msg":"n = 42","counters":{"error1":0,"error2":0},"stacktrace":"5:main.main [/source/go/scratch/scratch.go:23] < 6:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 7:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","baz":"waka","tags":["foo","bar","hoo"]}
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
func New(writer io.Writer, options *LoggerOptions) *Logger {

	logger := new(Logger)

	if options != nil {

		logger.options.AddSource = options.AddSource
		logger.options.BaseAttributes = options.BaseAttributes
		logger.options.BaseTags = options.BaseTags
		logger.options.Level = options.Level
		logger.options.ReplaceAttr = options.ReplaceAttr

	}

	if logger.options.Level == nil {
		logger.options.Level = new(slog.LevelVar)
	}

	hndlrOpts := new(slog.HandlerOptions)
	hndlrOpts.AddSource = logger.options.AddSource
	hndlrOpts.Level = logger.options.Level
	hndlrOpts.ReplaceAttr = newAttrReplacer(logger.options.ReplaceAttr)

	logger.wrapped = slog.New(slog.NewJSONHandler(writer, hndlrOpts))

	return logger
}

func (l *Logger) Trace(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, TRACE, message, attributes...)
}

func (l *Logger) Fine(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, FINE, message, attributes...)
}

func (l *Logger) Optional(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, OPTIONAL, message, attributes...)
}

func (l *Logger) Always(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, ALWAYS, message, attributes...)
}

func (l *Logger) OnPanic(ctx context.Context, handler PanicHandler, attributes ...any) {

	// note that there is no practical way to invoke this in automated unit
	// tests but see ../example/example.go for a demonstration

	if r := recover(); r != nil {

		msg, p := handler(r)

		l.Always(ctx, func() string { return msg }, attributes...)

		if p != nil {
			panic(p)
		}
	}
}

func (l *Logger) Enabled(ctx context.Context, verbosity Verbosity) bool {
	return l.wrapped.Enabled(ctx, slog.Level(verbosity))
}

func (l *Logger) BaseAttributes() []any {
	return l.options.BaseAttributes
}

func (l *Logger) SetBaseAttributes(attributes ...any) {
	l.options.BaseAttributes = attributes
}

func (l *Logger) BaseTags() []string {
	return l.options.BaseTags
}

func (l *Logger) SetBaseTags(tags ...string) {
	l.options.BaseTags = tags
}

func (l *Logger) Verbosity() Verbosity {
	return Verbosity(l.options.Level.Level())
}

func (l *Logger) SetVerbosity(verbosity Verbosity) {
	l.options.Level.Set(slog.Level(verbosity))
}
