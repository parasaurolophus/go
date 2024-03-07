// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"fmt"
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

	// Interface implemented by objects that wrap slog.Logger using a better
	// nomenclature for "levels."
	Logger interface {

		// Log at TRACE verbosity.
		Trace(message MessageBuilder, attributes ...any)

		// Log at TRACE verbosity using the supplied context.
		TraceContext(ctx context.Context, message MessageBuilder, attributes ...any)

		// Log at FINE verbosity.
		Fine(message MessageBuilder, attributes ...any)

		// Log at FINE verbosity using the supplied context.
		FineContext(ctx context.Context, message MessageBuilder, attributes ...any)

		// Log at OPTIONAL verbosity.
		Optional(message MessageBuilder, attributes ...any)

		// Log at OPTIONAL verbosity using the supplied context.
		OptionalContext(ctx context.Context, message MessageBuilder, attributes ...any)

		// Log at ALWAYS verbosity.
		Always(message MessageBuilder, attributes ...any)

		// Log at ALWAYS verbosity using the supplied context.
		AlwaysContext(ctx context.Context, message MessageBuilder, attributes ...any)
	}
)

func IsEnabled(logger Logger, verbosity Verbosity) bool {

	switch l := logger.(type) {

	case *syncLogger:
		return l.Enabled(verbosity)

	default:
		panic(fmt.Sprintf("can't set verbosity of loggers of type %T", l))
	}
}

func IsEnabledContext(logger Logger, ctx context.Context, verbosity Verbosity) bool {

	switch l := logger.(type) {

	case *syncLogger:
		return l.EnabledContext(ctx, verbosity)

	default:
		panic(fmt.Sprintf("can't set verbosity of loggers of type %T", l))
	}
}

func GetBaseAttributes(logger Logger) []any {

	switch l := logger.(type) {

	case *syncLogger:
		return l.BaseAttributes()

	default:
		panic(fmt.Sprintf("can't set base attributes of loggers of type %T", l))
	}
}

func SetBaseAttributes(logger Logger, attributes ...any) {

	switch l := logger.(type) {

	case *syncLogger:
		l.SetBaseAttributes(attributes...)

	default:
		panic(fmt.Sprintf("can't set base attributes of loggers of type %T", l))
	}
}

func GetBaseTags(logger Logger) []string {

	switch l := logger.(type) {

	case *syncLogger:
		return l.BaseTags()

	default:
		panic(fmt.Sprintf("can't set base tags of loggers of type %T", l))
	}
}

func SetBaseTags(logger Logger, tags ...string) {

	switch l := logger.(type) {

	case *syncLogger:
		l.SetBaseTags(tags...)

	default:
		panic(fmt.Sprintf("can't set base tags of loggers of type %T", l))
	}
}

func GetVerbosity(logger Logger) Verbosity {

	switch l := logger.(type) {

	case *syncLogger:
		return l.Verbosity()

	default:
		panic(fmt.Sprintf("can't set verbosity of loggers of type %T", l))
	}
}

func SetVerbosity(logger Logger, verbosity Verbosity) {

	switch l := logger.(type) {

	case *syncLogger:
		l.SetVerbosity(verbosity)

	default:
		panic(fmt.Sprintf("can't set verbosity of loggers of type %T", l))
	}
}
