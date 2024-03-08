// Copyright Kirk Rader 2024

package logging

import (
	"context"
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

		// Return true if and only if the specified verbosity is enabled.
		Enabled(verbosity Verbosity) bool

		// Return true if and only if the specified verbosity is enabled using
		// the specified context.
		EnabledContext(ctx context.Context, verbosity Verbosity) bool

		// Return the enabled verbosity level.
		Verbosity() Verbosity

		// Update the enabled verbosity level.
		SetVerbosity(verbosity Verbosity)

		// Return the base attributes.
		BaseAttributes() []any

		// Update the base attributes.
		SetBaseAttributes(attributes ...any)

		// Return the base tags.
		BaseTags() []string

		// Update the base tags.
		SetBaseTags(tags ...string)

		// Stop any asynchronous goroutines associated with this logger.
		Stop()
	}
)

const (

	// Value to use for FILE attribute when logging normal functionality.
	FILE_SKIPFRAMES_FOR_CALLER = -2

	// Value to use for FILE attribute when logging a panic.
	FILE_SKIPFRAMES_FOR_PANIC = -4
)
