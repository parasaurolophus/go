// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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

	// Type of function passed as first argument to Logger.Defer() and
	// Logger.DeferContext().
	Finally func()

	// Type of function passed to Logger.Defer() and Logger.DeferContext() to
	// allow for including the value returned by recover() in the log entry.
	RecoverHandler func(recovered any) string

	// Wrapper for an instance of slog.Logger.
	Logger struct {
		// This logger's configuration.
		options LoggerOptions
		// The wrapped slog.Logger.
		wrapped *slog.Logger
	}
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
func New(writer io.Writer, options *LoggerOptions) *Logger {

	logger := new(Logger)

	if options != nil {
		logger.options = *options
	}

	if logger.options.Level == nil {
		logger.options.Level = new(slog.LevelVar)
	}

	handlerOptions := new(slog.HandlerOptions)
	handlerOptions.AddSource = logger.options.AddSource
	handlerOptions.Level = logger.options.Level
	handlerOptions.ReplaceAttr = newAttrReplacer(logger.options.ReplaceAttr)
	logger.wrapped = slog.New(slog.NewJSONHandler(writer, handlerOptions))
	return logger
}

// Log at TRACE verbosity.
func (l *Logger) Trace(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, TRACE, message, attributes...)
}

// Log at TRACE verbosity using the supplied context.
func (l *Logger) TraceContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, TRACE, message, attributes...)
}

// Log at FINE verbosity.
func (l *Logger) Fine(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, FINE, message, attributes...)
}

// Log at FINE verbosity using the supplied context.
func (l *Logger) FineContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, FINE, message, attributes...)
}

// Log at OPTIONAL verbosity.
func (l *Logger) Optional(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, OPTIONAL, message, attributes...)
}

// Log at OPTIONAL verbosity using the supplied context.
func (l *Logger) OptionalContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, OPTIONAL, message, attributes...)
}

// Log at ALWAYS verbosity.
func (l *Logger) Always(message MessageBuilder, attributes ...any) {
	l.log(defaultContext, ALWAYS, message, attributes...)
}

// Log at ALWAYS verbosity using the supplied context.
func (l *Logger) AlwaysContext(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, ALWAYS, message, attributes...)
}

// See documentation for Logger.DeferContext().
func (l *Logger) Finally(panicAgain bool, finally Finally, recoverHandler RecoverHandler, attributes ...any) {
	l.finallyCommon(panicAgain, defaultContext, finally, recoverHandler, attributes...)
}

// For use with defer to log if a panic occurs.
//
// If recover() returns non-nil, its value will be passed to handler.
//
// Handler's return value will be used as the msg string in writing a log entry
// using l.AlwaysContext().
//
// If panicAgain is true, any panics that occur while this deferred method is in
// effect will be passed to panic() so as to cause the process to terminate
// abnormally.
//
// For example, if the following is invoked in a goroutine that was passed a
// channel named ch:
//
//	name := stacktraces.FunctionName()
//	defer logger.FinallyContext(
//
//	    // don't cause process to exit abnormally even if a panic occurs
//	    false,
//
//	    // clean-up function is always invoked
//	    func() { close(ch) },
//
//	    // remaining parameters are passed to logger.AlwaysContext() when
//	    // recover() returns non-nil
//
//	    ctx,
//	    func(r any) (string, any) {
//	        // second value will be used to resume panicing if non-nil
//	        // (typically this would be r to continue the now tidied
//	        // and logged panic in main.main or nil in a goroutine
//	        // so as to allow other goroutines to complete)
//	        return fmt.Sprintf("%s recovered from %v", name, r), nil
//	    },
//	)
//
// the goroutine will close ch on exit and, if a panic occurs, write a log entry
// whose msg is the string representation of the value returned by
// recover()while allowing other goroutines to continue running normally. If
// panicAgain were passed true, recovered value would be passed to panic() after
// the clean up and logging functions were invoked. The value of panicAgain is
// also used to determine whether or not panics in the clean-up or message
// builder functions cause an abnormal exit. [See the documentation for panic()
// and recover() for more information.]
func (l *Logger) FinallyContext(panicAgain bool, finally Finally, ctx context.Context, recoverHandler RecoverHandler, attributes ...any) {
	l.finallyCommon(panicAgain, ctx, finally, recoverHandler, attributes...)
}

// Return true or false depending on whether or not the given verbosity is
// currently enabled for the given logger.
func (l *Logger) Enabled(verbosity Verbosity) bool {
	return l.EnabledContext(defaultContext, verbosity)
}

// Return true or false depending on whether or not the given verbosity is
// currently enabled for the given logger.
func (l *Logger) EnabledContext(ctx context.Context, verbosity Verbosity) bool {
	return l.wrapped.Enabled(ctx, slog.Level(verbosity))
}

// Return the current base attributes.
func (l *Logger) BaseAttributes() []any {
	return l.options.BaseAttributes
}

// Update the base attributes.
func (l *Logger) SetBaseAttributes(attributes ...any) {
	l.options.BaseAttributes = attributes
}

// Return the current base tags.
func (l *Logger) BaseTags() []string {
	return l.options.BaseTags
}

// Update the base tags.
func (l *Logger) SetBaseTags(tags ...string) {
	l.options.BaseTags = tags
}

// Return the current verbosity
func (l *Logger) Verbosity() Verbosity {
	return Verbosity(l.options.Level.Level())
}

// Update the verbosity
func (l *Logger) SetVerbosity(verbosity Verbosity) {
	l.options.Level.Set(slog.Level(verbosity))
}

// Deprecated hack for backwards compatibility.
func (l *Logger) SetContext(ctx context.Context) {
	defaultContext = ctx
}

// Implement lazy evaluation of all log entry formatting code.
func (l *Logger) log(ctx context.Context, verbosity Verbosity, message MessageBuilder, attributes ...any) {

	if l.wrapped.Enabled(ctx, slog.Level(verbosity)) {

		msg := ""

		if message != nil {
			defer l.FinallyContext(l.options.AllowPanics, nil, ctx, nil)
			msg = message()
		}

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
				case TAGS:
					tags = appendTag(tags, combined[index+1])

				case STACKTRACE:
					stackTraceValue = convertSkipFrames(combined[index+1])
					includeStackTrace = true
				default:
					attribs = appendAttribute(attribs, attrib, combined[index+1])
				}
			}
		}

		if includeStackTrace {
			attribs = append(attribs, STACKTRACE, stacktraces.ShortStackTrace(stackTraceValue))
		}

		if len(tags) > 0 {
			attribs = append(attribs, TAGS, tags)
		}

		l.wrapped.Log(ctx, slog.Level(verbosity), msg, attribs...)
	}
}

// Return the result of appending the given key / val to the given attributes,
// so long as key is a string.
//
// Simply returns the given attributes list without appending anything if key is
// not a string.
func appendAttribute(attributes []any, key any, val any) []any {

	switch k := key.(type) {

	case string:
		return append(attributes, k, val)

	default:
		return attributes
	}
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

// Common implementation used by Logger.Finally() and Logger.FinallyContext()
func (l *Logger) finallyCommon(panicAgain bool, ctx context.Context, finally Finally, recoverHandler RecoverHandler, attributes ...any) {

	defer func() {
		if finally != nil {
			finally()
		}
	}()

	if r := recover(); r != nil {

		a := append(attributes, RECOVERED, r, STACKTRACE, -3)

		if recoverHandler == nil {
			l.AlwaysContext(ctx, nil, a...)
		} else {
			l.AlwaysContext(ctx, func() string { return recoverHandler(r) }, a...)
		}

		if panicAgain {
			panic(r)
		}
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