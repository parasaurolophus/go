// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"parasaurolophus/go/stacktraces"
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
		Trace(MessageBuilder, ...any)
		TraceContext(context.Context, MessageBuilder, ...any)
		Fine(MessageBuilder, ...any)
		FineContext(context.Context, MessageBuilder, ...any)
		Optional(MessageBuilder, ...any)
		OptionalContext(context.Context, MessageBuilder, ...any)
		Always(MessageBuilder, ...any)
		AlwaysContext(context.Context, MessageBuilder, ...any)
		Enabled(Verbosity) bool
		EnabledContext(context.Context, Verbosity) bool
		BaseAttributes() []any
		SetBaseAttributes(...any)
		BaseTags() []string
		SetBaseTags(...string)
		Verbosity() Verbosity
		SetVerbosity(Verbosity)
	}

	// Wrapper for an instance of slog.Logger.
	logger struct {
		// This logger's configuration.
		options LoggerOptions
		// The wrapped slog.Logger.
		wrapped *slog.Logger
	}
)

const (

	// Value to use for FILE attribute when logging normal functionality.
	FILE_SKIPFRAMES_FOR_CALLER = -2

	// Value to use for FILE attribute when logging a panic.
	FILE_SKIPFRAMES_FOR_PANIC = -4
)

var (

	// Default context for logging from command-line applications.
	defaultContext = context.Background()
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

	l := logger{}
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
	return &l
}

// Log at TRACE verbosity.
func (l *logger) Trace(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, TRACE, message, attributes...)
}

// Log at TRACE verbosity using the supplied context.
func (l *logger) TraceContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, TRACE, message, attributes...)
}

// Log at FINE verbosity.
func (l *logger) Fine(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, FINE, message, attributes...)
}

// Log at FINE verbosity using the supplied context.
func (l *logger) FineContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, FINE, message, attributes...)
}

// Log at OPTIONAL verbosity.
func (l *logger) Optional(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, OPTIONAL, message, attributes...)
}

// Log at OPTIONAL verbosity using the supplied context.
func (l *logger) OptionalContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, OPTIONAL, message, attributes...)
}

// Log at ALWAYS verbosity.
func (l *logger) Always(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, ALWAYS, message, attributes...)
}

// Log at ALWAYS verbosity using the supplied context.
func (l *logger) AlwaysContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, ALWAYS, message, attributes...)
}

// Return true or false depending on whether or not the given verbosity is
// currently enabled for the given logger.
func (l *logger) Enabled(verbosity Verbosity) bool {
	return l.EnabledContext(defaultContext, verbosity)
}

// Return true or false depending on whether or not the given verbosity is
// currently enabled for the given logger.
func (l *logger) EnabledContext(ctx context.Context, verbosity Verbosity) bool {
	return l.wrapped.Enabled(ctx, slog.Level(verbosity))
}

// Return the current base attributes.
func (l *logger) BaseAttributes() []any {
	return l.options.BaseAttributes
}

// Update the base attributes.
func (l *logger) SetBaseAttributes(attributes ...any) {
	l.options.BaseAttributes = attributes
}

// Return the current base tags.
func (l *logger) BaseTags() []string {
	return l.options.BaseTags
}

// Update the base tags.
func (l *logger) SetBaseTags(tags ...string) {
	l.options.BaseTags = tags
}

// Return the current verbosity
func (l *logger) Verbosity() Verbosity {
	return Verbosity(l.options.Level.Level())
}

// Update the verbosity
func (l *logger) SetVerbosity(verbosity Verbosity) {
	l.options.Level.Set(slog.Level(verbosity))
}

// Implement lazy evaluation of all log entry formatting code.
func (l *logger) log(context context.Context, verbosity Verbosity, messageBuilder MessageBuilder, attributes ...any) {

	if l.wrapped.Enabled(context, slog.Level(verbosity)) {

		msg := l.invokeMessageBuilder(context, messageBuilder)
		attribs := []any{}
		combined := append(l.options.BaseAttributes, attributes...)
		includeStackTrace := false
		var stackTraceValue any
		tags := l.options.BaseTags
		n := len(combined)
		var max int

		if n%2 == 0 {
			max = n - 1
		} else {
			max = n - 2
		}

		for index, attrib := range combined {
			if index < max && index%2 == 0 {
				switch attrib {

				case FILE:
					_, sourceInfo, ok := stacktraces.FunctionInfo(combined[index+1])
					if ok {
						attribs = l.appendAttribute(context, attribs, attrib, sourceInfo)
					} else {
						attribs = l.appendAttribute(context, attribs, attrib, combined[index+1])
						tags = appendTag(tags, FILE_ATTR_ERROR)
					}

				case STACKTRACE:
					stackTraceValue = convertSkipFrames(combined[index+1])
					includeStackTrace = true

				case TAGS:
					tags = appendTag(tags, combined[index+1])

				default:
					attribs = l.appendAttribute(context, attribs, attrib, combined[index+1])
				}
			}
		}

		if includeStackTrace {
			attribs = append(attribs, STACKTRACE, stacktraces.ShortStackTrace(stackTraceValue))
		}

		if len(tags) > 0 {
			attribs = append(attribs, TAGS, tags)
		}

		l.wrapped.Log(context, slog.Level(verbosity), msg, attribs...)
	}
}

// Safely invoke injected messageBuilder function.
func (l *logger) invokeMessageBuilder(context context.Context, messageBuilder MessageBuilder) string {

	defer l.logPanic(context)

	if messageBuilder != nil {
		return messageBuilder()
	}

	return ""
}

func (l *logger) logPanic(context context.Context) {

	if r := recover(); r != nil {
		l.AlwaysContext(
			context,
			func() string {
				return fmt.Sprintf("panic by message builder; recovered: %v", r)
			},
			TAGS, []string{PANIC, INJECTED},
			RECOVERED, r,
			STACKTRACE, nil,
		)
	}
}

// Return the result of appending the given key / val to the given attributes,
// so long as key is a string.
//
// Simply returns the given attributes list without appending anything if key is
// not a string.
func (l *logger) appendAttribute(context context.Context, attributes []any, key any, val any) []any {

	switch v := val.(type) {

	case func() any:
		val = l.invokeAttrHandler(context, v)
	}

	switch k := key.(type) {

	case string:
		return append(attributes, k, val)

	case fmt.Stringer:
		return append(attributes, k.String(), val)

	default:
		fmt.Fprintf(os.Stderr, "ignoring unsupported key, %v, of type %T", k, k)
		return attributes
	}
}

// Safely invoke injected attribute value function.
func (l *logger) invokeAttrHandler(context context.Context, handler func() any) any {

	defer l.logPanic(context)
	return handler()
}

// Return the result of merging base tags with those supplied as the value of a
// "tags" attribute in a call to a logging method.
//
// This is very permissive in interpreting the tag to be appended. All the
// elements in a slice of strings will be appended. A single string will be
// appended. Any other type will be converted to a string and appended.
func appendTag(tags []string, tag any) []string {

	switch v := tag.(type) {

	case []string:
		return append(tags, v...)

	case string:
		return append(tags, v)

	case fmt.Stringer:
		return append(tags, v.String())

	default:
		return append(tags, fmt.Sprintf("%v", v))
	}
}

// Each time the function returned by this one is called it will:
//
//   - Invoke oldReplacer if it is not nil.
//
//   - Replace "level" with "verbosity" whose value is the name of the
//     corresponding Verbosity const ("TRACE", "FINE", "OPTIONAL" or "ALWAYS").
func newAttrReplacer(oldReplacer func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {

	return func(groups []string, attr slog.Attr) slog.Attr {

		if oldReplacer != nil {
			attr = oldReplacer(groups, attr)
		}

		if attr.Key == "level" {

			const verbosityKey = "verbosity"
			val := attr.Value.String()

			switch val {

			case slog.LevelDebug.String():
				return slog.Attr{
					Key:   verbosityKey,
					Value: slog.StringValue(TRACE.String()),
				}

			case slog.LevelInfo.String():
				return slog.Attr{
					Key:   verbosityKey,
					Value: slog.StringValue(FINE.String()),
				}

			case slog.LevelWarn.String():
				return slog.Attr{
					Key:   verbosityKey,
					Value: slog.StringValue(OPTIONAL.String()),
				}

			case slog.LevelError.String():
				return slog.Attr{
					Key:   verbosityKey,
					Value: slog.StringValue(ALWAYS.String()),
				}

			default:
				return slog.Attr{
					Key:   verbosityKey,
					Value: slog.StringValue(val),
				}
			}
		}

		return attr
	}
}

// Return the value to use as the skipFrames parameter for the value of a given
// "stacktrace" attribute.
func convertSkipFrames(val any) any {

	// the number of frames to skip here is empirically
	// derived and may change if this library is refactored
	const defaultSkip = 5

	switch v := val.(type) {

	case int:
		if v < 0 {
			return defaultSkip - v
		}
		return v

	case string:
		return v

	default:
		return defaultSkip
	}
}
