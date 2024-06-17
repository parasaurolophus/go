// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"parasaurolophus/go/stacktraces"
)

// Wrapper for an instance of slog.Logger.
type loggerStruct struct {
	// This logger's configuration.
	options LoggerOptions
	// The wrapped slog.Logger.
	wrapped *slog.Logger
}

// Log at TRACE verbosity using the supplied context.
func (l *loggerStruct) Trace(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, TRACE, message, attributes...)
}

// Log at FINE verbosity using the supplied context.
func (l *loggerStruct) Fine(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, FINE, message, attributes...)
}

// Log at OPTIONAL verbosity using the supplied context.
func (l *loggerStruct) Optional(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, OPTIONAL, message, attributes...)
}

// Log at ALWAYS verbosity using the supplied context.
func (l *loggerStruct) Always(ctx context.Context, message MessageBuilder, attributes ...any) {
	l.log(ctx, ALWAYS, message, attributes...)
}

// Return true or false depending on whether or not the given verbosity is
// currently enabled for the given logger.
func (l *loggerStruct) Enabled(ctx context.Context, verbosity Verbosity) bool {
	return l.wrapped.Enabled(ctx, slog.Level(verbosity))
}

// Return the current base attributes.
func (l *loggerStruct) BaseAttributes() []any {
	return l.options.BaseAttributes
}

// Update the base attributes.
func (l *loggerStruct) SetBaseAttributes(attributes ...any) {
	l.options.BaseAttributes = attributes
}

// Return the current base tags.
func (l *loggerStruct) BaseTags() []string {
	return l.options.BaseTags
}

// Update the base tags.
func (l *loggerStruct) SetBaseTags(tags ...string) {
	l.options.BaseTags = tags
}

// Return the current verbosity
func (l *loggerStruct) Verbosity() Verbosity {
	return Verbosity(l.options.Level.Level())
}

// Update the verbosity
func (l *loggerStruct) SetVerbosity(verbosity Verbosity) {
	l.options.Level.Set(slog.Level(verbosity))
}

// Implement lazy evaluation of all log entry formatting code.
func (l *loggerStruct) log(context context.Context, verbosity Verbosity, messageBuilder MessageBuilder, attributes ...any) {

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
func (l *loggerStruct) invokeMessageBuilder(context context.Context, messageBuilder MessageBuilder) string {

	defer l.logPanic(context)

	if messageBuilder != nil {
		return messageBuilder()
	}

	return ""
}

func (l *loggerStruct) logPanic(context context.Context) {

	if r := recover(); r != nil {
		l.Always(
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
func (l *loggerStruct) appendAttribute(context context.Context, attributes []any, key any, val any) []any {

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
func (l *loggerStruct) invokeAttrHandler(context context.Context, handler func() any) any {

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
