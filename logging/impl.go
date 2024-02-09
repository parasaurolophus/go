package logging

import (
	"context"
	"fmt"
	"log/slog"
	"parasaurolophus/go/stacktraces"
)

// Implement lazy evaluation of all log entry formatting code.
//
// Only invoke message builder if the wrapped slog.Logger is enabled for the
// slog.Level equivalent of the specified Verbosity.
//
// Handles base attributes by combining them with any that are explicitly
// supplied as parameters to this method.
//
// Handles base tags by combining them with the value of any explictly provided
// "tags" attribute or adding a "tags" attribute when none was supplied.
//
// Handles a "stacktrace" attribute by replacing its value with a one-line stack
// trace starting at the point in the call stack just beyond the stacktraces and
// logging library internals.
func (l *Logger) log(ctx context.Context, verbosity Verbosity, message MessageBuilder, attributes ...any) {

	if l.wrapped.Enabled(ctx, slog.Level(verbosity)) {

		msg := ""

		if message != nil {
			msg = message()
		}

		attribs := []any{}
		combined := append(l.options.BaseAttributes, attributes...)
		tags := l.options.BaseTags
		max := len(combined) - 1

		for index, attrib := range combined {

			if index < max && index%2 == 0 {

				if attrib == TAGS {

					switch v := combined[index+1].(type) {

					case []string:
						tags = append(tags, v...)

					case string:
						tags = append(tags, v)

					default:
						tags = append(tags, fmt.Sprintf("%v", v))
					}

				} else if attrib == STACKTRACE {

					// the number of stack frames to skip is empirically derived
					// and may change as a result of any code refactoring within
					// the logging and stacktraces packages
					attribs = append(attribs, STACKTRACE, stacktraces.ShortStackTrace(5))

				} else {

					attribs = append(attribs, attrib, combined[index+1])
				}
			}
		}

		if len(tags) > 0 {

			attribs = append(attribs, TAGS, tags)
		}

		l.wrapped.Log(ctx, slog.Level(verbosity), msg, attribs...)
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

			switch attr.Value.String() {

			case "DEBUG":
				return slog.Attr{
					Key:   "verbosity",
					Value: slog.StringValue("TRACE"),
				}

			case "INFO":
				return slog.Attr{
					Key:   "verbosity",
					Value: slog.StringValue("FINE"),
				}

			case "WARN":
				return slog.Attr{
					Key:   "verbosity",
					Value: slog.StringValue("OPTIONAL"),
				}

			case "ERROR":
				return slog.Attr{
					Key:   "verbosity",
					Value: slog.StringValue("ALWAYS"),
				}

			default:
				return slog.Attr{
					Key:   "verbosity",
					Value: slog.StringValue(fmt.Sprint(attr.Value.String())),
				}
			}
		}

		return attr
	}
}
